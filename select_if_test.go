package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectBuilderConditional(t *testing.T) {
	baseQuery := Select("id").From("users")

	t.Run("DistinctIf", func(t *testing.T) {
		b := Select("id").From("users").DistinctIf(true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT DISTINCT id FROM users", sql)

		b = Select("id").From("users").DistinctIf(false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("ColumnIf", func(t *testing.T) {
		b := baseQuery.ColumnIf(Expr("COUNT(*)"), true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, COUNT(*) FROM users", sql)

		b = baseQuery.ColumnIf(Expr("COUNT(*)"), false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("ColumnsIf (plural)", func(t *testing.T) {
		b := Select().From("users").ColumnsIf(
			ValIf[safeString]("id", true),
			ValIf[safeString]("deleted_at", false),
			ValIf[safeString]("name", true),
		)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id, name FROM users", sql)
	})

	t.Run("JoinIf", func(t *testing.T) {
		b := baseQuery.JoinIf("posts ON posts.user_id = users.id", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users JOIN posts ON posts.user_id = users.id", sql)

		b = baseQuery.LeftJoinIf("posts ON posts.user_id = users.id", true)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users LEFT JOIN posts ON posts.user_id = users.id", sql)

		b = baseQuery.RightJoinIf("posts ON posts.user_id = users.id", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("WhereIf", func(t *testing.T) {
		b := baseQuery.WhereIf(Eq{"status": "active"}, true)
		sql, args, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users WHERE status = ?", sql)
		assert.Equal(t, []interface{}{"active"}, args)

		b = baseQuery.WhereIf(Eq{"status": "active"}, false)
		sql, args, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
		assert.Empty(t, args)
	})

	t.Run("GroupByIf (singular)", func(t *testing.T) {
		b := baseQuery.GroupByIf("status", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users GROUP BY status", sql)

		b = baseQuery.GroupByIf("status", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("GroupBysIf (plural)", func(t *testing.T) {
		b := baseQuery.GroupBysIf(
			ValIf[safeString]("status", true),
			ValIf[safeString]("company_id", false),
			ValIf[safeString]("department", true),
		)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users GROUP BY status, department", sql)
	})

	t.Run("HavingIf", func(t *testing.T) {
		b := baseQuery.GroupBy("status").HavingIf(Expr("COUNT(*) > ?", 10), true)
		sql, args, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users GROUP BY status HAVING COUNT(*) > ?", sql)
		assert.Equal(t, []interface{}{10}, args)

		b = baseQuery.GroupBy("status").HavingIf(Expr("COUNT(*) > ?", 10), false)
		sql, args, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users GROUP BY status", sql)
		assert.Empty(t, args)
	})

	t.Run("OrderByIf (singular)", func(t *testing.T) {
		b := baseQuery.OrderByIf("created_at DESC", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users ORDER BY created_at DESC", sql)

		b = baseQuery.OrderByIf("created_at DESC", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("LimitIf", func(t *testing.T) {
		b := baseQuery.LimitIf(50, true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users LIMIT 50", sql)

		b = baseQuery.LimitIf(50, false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("OffsetIf", func(t *testing.T) {
		b := baseQuery.OffsetIf(100, true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users OFFSET 100", sql)

		b = baseQuery.OffsetIf(100, false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})

	t.Run("PrefixIf & SuffixIf", func(t *testing.T) {
		b := baseQuery.PrefixIf("HINT_A", true).SuffixIf("HINT_B", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "HINT_A SELECT id FROM users HINT_B", sql)

		b = baseQuery.PrefixIf("HINT_A", false).SuffixIf("HINT_B", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT id FROM users", sql)
	})
}
