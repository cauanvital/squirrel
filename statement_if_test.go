package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectIfFunctions(t *testing.T) {
	t.Run("Package Level SelectIf", func(t *testing.T) {

		t.Run("should select all columns when all are included", func(t *testing.T) {
			b := SelectIf(
				ValIf[safeString]("id", true),
				ValIf[safeString]("name", true),
				ValIf[safeString]("email", true),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT id, name, email FROM users", sql)
		})

		t.Run("should select only included columns from a mixed set", func(t *testing.T) {
			b := SelectIf(
				ValIf[safeString]("id", true),
				ValIf[safeString]("password_hash", false),
				ValIf[safeString]("email", true),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT id, email FROM users", sql)
		})

		t.Run("should default to SELECT * when all columns are excluded", func(t *testing.T) {
			b := SelectIf(
				ValIf[safeString]("id", false),
				ValIf[safeString]("name", false),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT * FROM users", sql)
		})

		t.Run("should default to SELECT * when called with no arguments", func(t *testing.T) {
			b := SelectIf().From("users")
			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT * FROM users", sql)
		})
	})

	t.Run("StatementBuilderType Method SelectIf", func(t *testing.T) {

		t.Run("should select all columns when all are included", func(t *testing.T) {
			b := StatementBuilder.SelectIf(
				ValIf[safeString]("id", true),
				ValIf[safeString]("name", true),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT id, name FROM users", sql)
		})

		t.Run("should select only included columns from a mixed set", func(t *testing.T) {
			b := StatementBuilder.SelectIf(
				ValIf[safeString]("id", true),
				ValIf[safeString]("last_login", false),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT id FROM users", sql)
		})

		t.Run("should default to SELECT * when all columns are excluded", func(t *testing.T) {
			b := StatementBuilder.SelectIf(
				ValIf[safeString]("id", false),
			).From("users")

			sql, _, err := b.ToSql()
			assert.NoError(t, err)
			assert.Equal(t, "SELECT * FROM users", sql)
		})
	})
}
