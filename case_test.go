package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseBuilder(t *testing.T) {
	t.Run("Immutability is preserved", func(t *testing.T) {
		b1 := CaseBuilder()
		b2 := b1.When(Eq{"a": 1}, Expr("100"))
		b3 := b2.Else(Expr("0"))

		_, _, err1 := b1.ToSql()
		assert.Error(t, err1, "Original builder should be unchanged and invalid")

		sql2, args2, err2 := b2.ToSql()
		assert.NoError(t, err2)
		assert.Equal(t, "CASE WHEN a = ? THEN 100 END", sql2)
		assert.Equal(t, []interface{}{1}, args2)

		sql3, args3, err3 := b3.ToSql()
		assert.NoError(t, err3)
		assert.Equal(t, "CASE WHEN a = ? THEN 100 ELSE 0 END", sql3)
		assert.Equal(t, []interface{}{1}, args3)
	})

	t.Run("Builds simple CASE with one WHEN", func(t *testing.T) {
		b := CaseBuilder().When(Expr("status = ?", "active"), Expr("'active_user'"))
		sql, args, err := b.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, "CASE WHEN status = ? THEN 'active_user' END", sql)
		assert.Equal(t, []interface{}{"active"}, args)
	})

	t.Run("Builds CASE with multiple WHENs", func(t *testing.T) {
		b := CaseBuilder().
			When(Eq{"color": "red"}, Expr("1")).
			When(Eq{"color": "blue"}, Expr("2"))

		sql, args, err := b.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, "CASE WHEN color = ? THEN 1 WHEN color = ? THEN 2 END", sql)
		assert.Equal(t, []interface{}{"red", "blue"}, args)
	})

	t.Run("Builds CASE with an ELSE part", func(t *testing.T) {
		b := CaseBuilder().
			When(Gt{"age": 18}, Expr("'adult'")).
			Else(Expr("'minor'"))

		sql, args, err := b.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, "CASE WHEN age > ? THEN 'adult' ELSE 'minor' END", sql)
		assert.Equal(t, []interface{}{18}, args)
	})

	t.Run("Builds CASE with multiple WHENs and an ELSE part", func(t *testing.T) {
		b := CaseBuilder().
			When(Eq{"category": "A"}, Expr("10")).
			When(Eq{"category": "B"}, Expr("20")).
			Else(Expr("30"))

		sql, args, err := b.ToSql()

		assert.NoError(t, err)
		assert.Equal(t, "CASE WHEN category = ? THEN 10 WHEN category = ? THEN 20 ELSE 30 END", sql)
		assert.Equal(t, []interface{}{"A", "B"}, args)
	})

	t.Run("ToSql returns error if no WHEN clauses are set", func(t *testing.T) {
		_, _, err := CaseBuilder().ToSql()
		assert.Error(t, err)
		assert.Equal(t, "case expression must contain at lease one WHEN clause", err.Error())
	})

	t.Run("MustSql panics if no WHEN clauses are set", func(t *testing.T) {
		panicked := false
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			assert.True(t, panicked, "MustSql should panic when ToSql would return an error")
		}()

		CaseBuilder().MustSql()
	})

	t.Run("Integration with SelectBuilder", func(t *testing.T) {
		caseStmt := CaseBuilder().
			When(Eq{"status": 1}, Expr("'active'")).
			When(Eq{"status": 2}, Expr("'pending'")).
			Else(Expr("'inactive'"))

		b := SelectBuilder().
			Column(Alias(caseStmt, "user_status")).
			From("users")

		sql, args, err := b.ToSql()
		assert.NoError(t, err)

		expectedSql := "SELECT (CASE WHEN status = ? THEN 'active' WHEN status = ? THEN 'pending' ELSE 'inactive' END) AS user_status FROM users"
		assert.Equal(t, expectedSql, sql)
		assert.Equal(t, []interface{}{1, 2}, args)
	})
}
