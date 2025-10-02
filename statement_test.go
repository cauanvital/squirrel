package squirrel2

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatementBuilder(t *testing.T) {
	db := &DBStub{}
	sb := StatementBuilder.RunWith(db)

	sb.Select("test").Exec()
	assert.Equal(t, "SELECT test", db.LastExecSql)
}

func TestStatementBuilderPlaceholderFormat(t *testing.T) {
	db := &DBStub{}
	sb := StatementBuilder.RunWith(db).PlaceholderFormat(Dollar)

	sb.Select("test").Where(Expr("x = ?")).Exec()
	assert.Equal(t, "SELECT test WHERE x = $1", db.LastExecSql)
}

func TestRunWithDB(t *testing.T) {
	db := &sql.DB{}
	assert.NotPanics(t, func() {
		Select().RunWith(db)
		Insert("t").RunWith(db)
		Update("t").RunWith(db)
		Delete("t").RunWith(db)
	}, "RunWith(*sql.DB) should not panic")

}

func TestRunWithTx(t *testing.T) {
	tx := &sql.Tx{}
	assert.NotPanics(t, func() {
		Select().RunWith(tx)
		Insert("t").RunWith(tx)
		Update("t").RunWith(tx)
		Delete("t").RunWith(tx)
	}, "RunWith(*sql.Tx) should not panic")
}

type fakeBaseRunner struct{}

func (fakeBaseRunner) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (fakeBaseRunner) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func TestRunWithBaseRunner(t *testing.T) {
	sb := StatementBuilder.RunWith(fakeBaseRunner{})
	_, err := sb.Select("test").Exec()
	assert.NoError(t, err)
}

func TestRunWithBaseRunnerQueryRowError(t *testing.T) {
	sb := StatementBuilder.RunWith(fakeBaseRunner{})
	assert.Error(t, ErrRunnerNotQueryRunner, sb.Select("test").QueryRow().Scan(nil))

}

func TestStatementBuilderWhere(t *testing.T) {
	sb := StatementBuilder.Where(Expr("x = ?", 1))

	sql, args, err := sb.Select("test").Where(Expr("y = ?", 2)).ToSql()
	assert.NoError(t, err)

	expectedSql := "SELECT test WHERE x = ? AND y = ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{1, 2}
	assert.Equal(t, expectedArgs, args)
}
