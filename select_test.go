package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("ToSql builds a complete query", func(t *testing.T) {
		b := SelectBuilder(StatementBuilder).
			Prefix("WITH my_cte AS (SELECT 1)").
			Distinct().
			Columns("id", "name").
			Column(Expr("COUNT(*) AS total")).
			From("users u").
			Join("posts p ON p.user_id = u.id").
			LeftJoin("comments c ON c.post_id = p.id").
			Where(Eq{"u.status": "active"}).
			GroupBy("u.id", "u.name").
			Having(Expr("COUNT(*) > ?", 1)).
			OrderBy("u.name DESC").
			Limit(10).
			Offset(20).
			Suffix("FETCH NEXT 10 ROWS ONLY")

		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "WITH my_cte AS (SELECT 1) " +
			"SELECT DISTINCT id, name, COUNT(*) AS total " +
			"FROM users u " +
			"JOIN posts p ON p.user_id = u.id " +
			"LEFT JOIN comments c ON c.post_id = p.id " +
			"WHERE u.status = ? " +
			"GROUP BY u.id, u.name " +
			"HAVING COUNT(*) > ? " +
			"ORDER BY u.name DESC " +
			"LIMIT 10 OFFSET 20 " +
			"FETCH NEXT 10 ROWS ONLY"
		assert.Equal(t, expectedSql, sql)

		expectedArgs := []interface{}{"active", 1}
		assert.Equal(t, expectedArgs, args)
	})

	t.Run("Immutability is preserved", func(t *testing.T) {
		// This is the most important test for your new design.
		// It ensures that the builder is truly immutable.
		b1 := SelectBuilder(StatementBuilder).Columns("id").From("users")
		b2 := b1.Where(Eq{"status": "active"})
		b3 := b2.Limit(5)

		// Get SQL from all three builders
		sql1, args1, err1 := b1.ToSql()
		assert.NoError(t, err1)

		sql2, args2, err2 := b2.ToSql()
		assert.NoError(t, err2)

		sql3, args3, err3 := b3.ToSql()
		assert.NoError(t, err3)

		// Assert that the original builders were not modified
		assert.Equal(t, "SELECT id FROM users", sql1)
		assert.Empty(t, args1)

		assert.Equal(t, "SELECT id FROM users WHERE status = ?", sql2)
		assert.Equal(t, []interface{}{"active"}, args2)

		assert.Equal(t, "SELECT id FROM users WHERE status = ? LIMIT 5", sql3)
		assert.Equal(t, []interface{}{"active"}, args3)

		// Also assert that they are not equal to each other
		assert.NotEqual(t, sql1, sql2)
		assert.NotEqual(t, sql2, sql3)
	})

	t.Run("ToSql returns error if no columns are set", func(t *testing.T) {
		_, _, err := SelectBuilder(StatementBuilder).From("users").ToSql()
		assert.Error(t, err)
		assert.Equal(t, "select statements must have at least one result column", err.Error())
	})

	t.Run("RemoveColumns works correctly", func(t *testing.T) {
		b1 := SelectBuilder(StatementBuilder).Columns("id", "name").From("users")
		b2 := b1.RemoveColumns()
		b3 := b2.Column(safeString("email"))

		// The original builder is untouched
		sql1, _, _ := b1.ToSql()
		assert.Equal(t, "SELECT id, name FROM users", sql1)

		// The builder with removed columns should error
		_, _, err := b2.ToSql()
		assert.Error(t, err, "builder with removed columns should fail ToSql")

		// The builder with a new column should work
		sql3, _, _ := b3.ToSql()
		assert.Equal(t, "SELECT email FROM users", sql3)
	})

	t.Run("RemoveLimit and RemoveOffset work correctly", func(t *testing.T) {
		b1 := SelectBuilder(StatementBuilder).Columns("*").From("t").Limit(10).Offset(20)
		b2 := b1.RemoveLimit()
		b3 := b2.RemoveOffset()

		sql1, _, _ := b1.ToSql()
		assert.Equal(t, "SELECT * FROM t LIMIT 10 OFFSET 20", sql1)

		sql2, _, _ := b2.ToSql()
		assert.Equal(t, "SELECT * FROM t OFFSET 20", sql2)

		sql3, _, _ := b3.ToSql()
		assert.Equal(t, "SELECT * FROM t", sql3)
	})

	t.Run("PlaceholderFormat changes placeholders", func(t *testing.T) {
		b := SelectBuilder(StatementBuilder).
			Columns("id").
			From("users").
			Where(Eq{"id": 1, "name": "foo"}).
			PlaceholderFormat(Dollar)

		sql, _, err := b.ToSql()
		assert.NoError(t, err)

		// The order of columns in Eq is not guaranteed, so we check both possibilities
		possibleSql1 := "SELECT id FROM users WHERE id = $1 AND name = $2"
		possibleSql2 := "SELECT id FROM users WHERE name = $1 AND id = $2"
		assert.Contains(t, []string{possibleSql1, possibleSql2}, sql)
	})

	t.Run("FromSelect builds subquery correctly", func(t *testing.T) {
		subquery := SelectBuilder(StatementBuilder).
			Columns("id").
			From("posts").
			Where(Eq{"published": true})

		b := SelectBuilder(StatementBuilder).
			Columns("user_id").
			FromSelect(subquery, "p").
			GroupBy("user_id")

		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "SELECT user_id FROM (SELECT id FROM posts WHERE published = ?) AS p GROUP BY user_id"
		assert.Equal(t, expectedSql, sql)
		assert.Equal(t, []interface{}{true}, args)
	})
}
