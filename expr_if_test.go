package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalSqlizers(t *testing.T) {
	t.Run("ExprIf", func(t *testing.T) {
		t.Run("should be excluded when include is false", func(t *testing.T) {
			e := ExprIf(Expr("a = ?", 1), false)

			sql, args, err := e.ToSql()

			assert.NoError(t, err)
			assert.Empty(t, sql, "SQL should be empty when include is false")
			assert.Nil(t, args, "Args should be nil when include is false")
		})

		t.Run("should be included with simple arguments when include is true", func(t *testing.T) {
			e := ExprIf(Expr("a = ? AND b = ?", 1, "two"), true)

			sql, args, err := e.ToSql()

			assert.NoError(t, err)
			assert.Equal(t, "a = ? AND b = ?", sql)
			assert.Equal(t, []interface{}{1, "two"}, args)
		})

		t.Run("should expand nested Sqlizer when include is true", func(t *testing.T) {
			nestedClause := Eq{"status": "active"}
			e := ExprIf(Expr("user_id = ? AND details IN (?)", 123, nestedClause), true)

			sql, args, err := e.ToSql()

			assert.NoError(t, err)
			assert.Equal(t, "user_id = ? AND details IN (status = ?)", sql)
			assert.Equal(t, []interface{}{123, "active"}, args)
		})
	})
}
