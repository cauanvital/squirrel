package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("ToSql builds a complete query", func(t *testing.T) {
		b := DeleteBuilder().
			Prefix("WITH my_cte AS (SELECT 1)").
			From("users").
			Where(Eq{"id": 1}).
			Where(Lt{"age": 30}).
			OrderBy("name DESC").
			Limit(5).
			Offset(10).
			Suffix("RETURNING id")

		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "WITH my_cte AS (SELECT 1) " +
			"DELETE FROM users " +
			"WHERE id = ? AND age < ? " +
			"ORDER BY name DESC " +
			"LIMIT 5 OFFSET 10 " +
			"RETURNING id"
		assert.Equal(t, expectedSql, sql)

		expectedArgs := []interface{}{1, 30}
		assert.Equal(t, expectedArgs, args)
	})

	t.Run("Immutability is preserved", func(t *testing.T) {
		b1 := DeleteBuilder().From("users")
		b2 := b1.Where(Eq{"status": "archived"})
		b3 := b2.Limit(100)

		sql1, args1, err1 := b1.ToSql()
		assert.NoError(t, err1)

		sql2, args2, err2 := b2.ToSql()
		assert.NoError(t, err2)

		sql3, args3, err3 := b3.ToSql()
		assert.NoError(t, err3)

		assert.Equal(t, "DELETE FROM users", sql1)
		assert.Empty(t, args1)

		assert.Equal(t, "DELETE FROM users WHERE status = ?", sql2)
		assert.Equal(t, []interface{}{"archived"}, args2)

		assert.Equal(t, "DELETE FROM users WHERE status = ? LIMIT 100", sql3)
		assert.Equal(t, []interface{}{"archived"}, args3)

		assert.NotEqual(t, sql1, sql2)
		assert.NotEqual(t, sql2, sql3)
	})

	t.Run("ToSql returns error if From is not set", func(t *testing.T) {
		_, _, err := DeleteBuilder().Where(Eq{"id": 1}).ToSql()
		assert.Error(t, err)
		assert.Equal(t, "delete statements must specify a From table", err.Error())
	})

	t.Run("ToSql builds correctly without a Where clause", func(t *testing.T) {
		sql, args, err := DeleteBuilder().From("logs").ToSql()
		assert.NoError(t, err)

		assert.Equal(t, "DELETE FROM logs", sql)
		assert.Empty(t, args)
	})

	t.Run("PlaceholderFormat changes placeholders", func(t *testing.T) {
		b := DeleteBuilder().
			From("users").
			Where(Eq{"id": 1}).
			PlaceholderFormat(Dollar)

		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, "DELETE FROM users WHERE id = $1", sql)
		assert.Equal(t, []interface{}{1}, args)
	})

	t.Run("MustSql panics on error", func(t *testing.T) {
		panicked := false
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			assert.True(t, panicked, "MustSql should panic on error")
		}()

		DeleteBuilder().Where(Eq{"id": 1}).MustSql()
	})
}
