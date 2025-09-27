package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilderConditional(t *testing.T) {
	baseBuilder := DeleteBuilder().From("users")

	t.Run("Immutability is preserved", func(t *testing.T) {
		b1 := DeleteBuilder().From("users")
		b2 := b1.WhereIf(Eq{"id": 1}, true)
		b3 := b2.LimitIf(5, true)

		sql1, _, _ := b1.ToSql()
		assert.Equal(t, "DELETE FROM users", sql1)

		sql2, _, _ := b2.ToSql()
		assert.Equal(t, "DELETE FROM users WHERE id = ?", sql2)

		sql3, _, _ := b3.ToSql()
		assert.Equal(t, "DELETE FROM users WHERE id = ? LIMIT 5", sql3)

		assert.NotEqual(t, sql1, sql2)
		assert.NotEqual(t, sql2, sql3)
	})

	t.Run("WhereIf", func(t *testing.T) {
		b := baseBuilder.WhereIf(Eq{"status": "banned"}, true)
		sql, args, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users WHERE status = ?", sql)
		assert.Equal(t, []interface{}{"banned"}, args)

		b = baseBuilder.WhereIf(Eq{"status": "banned"}, false)
		sql, args, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users", sql)
		assert.Empty(t, args)
	})

	t.Run("LimitIf", func(t *testing.T) {
		b := baseBuilder.LimitIf(10, true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users LIMIT 10", sql)

		b = baseBuilder.LimitIf(10, false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users", sql)
	})

	t.Run("OffsetIf", func(t *testing.T) {
		b := baseBuilder.OffsetIf(20, true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users OFFSET 20", sql)

		b = baseBuilder.OffsetIf(20, false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users", sql)
	})

	t.Run("OrderByIf (singular)", func(t *testing.T) {
		b := baseBuilder.OrderByIf("created_at DESC", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users ORDER BY created_at DESC", sql)

		b = baseBuilder.OrderByIf("created_at DESC", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users", sql)
	})

	t.Run("OrderBysIf (plural)", func(t *testing.T) {
		b := baseBuilder.OrderBysIf(
			valIf[safeString]{Value: "name ASC", Include: true},
			valIf[safeString]{Value: "email ASC", Include: false},
			valIf[safeString]{Value: "id DESC", Include: true},
		)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users ORDER BY name ASC, id DESC", sql)
	})

	t.Run("PrefixIf and SuffixIf", func(t *testing.T) {
		b := baseBuilder.PrefixIf("/* trace_id=123 */", true).SuffixIf("RETURNING id", true)
		sql, _, err := b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "/* trace_id=123 */ DELETE FROM users RETURNING id", sql)

		b = baseBuilder.PrefixIf("/* trace_id=123 */", true).SuffixIf("RETURNING id", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "/* trace_id=123 */ DELETE FROM users", sql)

		b = baseBuilder.PrefixIf("/* trace_id=123 */", false).SuffixIf("RETURNING id", false)
		sql, _, err = b.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users", sql)
	})
}
