package squirrel

import (
	"database/sql"
	"sort"

	"github.com/lann/builder"
)

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b InsertBuilder) PlaceholderFormat(f PlaceholderFormat) InsertBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(InsertBuilder)
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b InsertBuilder) RunWith(runner BaseRunner) InsertBuilder {
	return setRunWith(b, runner).(InsertBuilder)
}

// Exec builds and Execs the query with the Runner set by RunWith.
func (b InsertBuilder) Exec() (sql.Result, error) {
	data := builder.GetStruct(b).(insertData)
	return data.Exec()
}

// Query builds and Querys the query with the Runner set by RunWith.
func (b InsertBuilder) Query() (*sql.Rows, error) {
	data := builder.GetStruct(b).(insertData)
	return data.Query()
}

// QueryRow builds and QueryRows the query with the Runner set by RunWith.
func (b InsertBuilder) QueryRow() RowScanner {
	data := builder.GetStruct(b).(insertData)
	return data.QueryRow()
}

// Scan is a shortcut for QueryRow().Scan.
func (b InsertBuilder) Scan(dest ...interface{}) error {
	return b.QueryRow().Scan(dest...)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b InsertBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(insertData)
	return data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b InsertBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query
func (b InsertBuilder) Prefix(sql safeString, args ...interface{}) InsertBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b InsertBuilder) PrefixExpr(expr Sqlizer) InsertBuilder {
	return builder.Append(b, "Prefixes", expr).(InsertBuilder)
}

// Options adds keyword options before the INTO clause of the query.
func (b InsertBuilder) Options(options ...safeString) InsertBuilder {
	return builder.Extend(b, "Options", options).(InsertBuilder)
}

// Into sets the INTO clause of the query.
func (b InsertBuilder) Into(into safeString) InsertBuilder {
	return builder.Set(b, "Into", into).(InsertBuilder)
}

// Columns adds insert columns to the query.
func (b InsertBuilder) Columns(columns ...safeString) InsertBuilder {
	return builder.Extend(b, "Columns", columns).(InsertBuilder)
}

// Values adds a single row's values to the query.
func (b InsertBuilder) Values(values ...interface{}) InsertBuilder {
	return builder.Append(b, "Values", values).(InsertBuilder)
}

// Suffix adds an expression to the end of the query
func (b InsertBuilder) Suffix(sql safeString, args ...interface{}) InsertBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b InsertBuilder) SuffixExpr(expr Sqlizer) InsertBuilder {
	return builder.Append(b, "Suffixes", expr).(InsertBuilder)
}

// SetMap set columns and values for insert builder from a map of column name and value
// note that it will reset all previous columns and values was set if any
func (b InsertBuilder) SetMap(clauses map[safeString]interface{}) InsertBuilder {
	// Keep the columns in a consistent order by sorting the column key string.
	cols := make([]safeString, 0, len(clauses))
	for col := range clauses {
		cols = append(cols, col)
	}
	sort.Slice(cols, func(idx, jdx int) bool {
		return cols[idx] < cols[jdx]
	})

	vals := make([]interface{}, 0, len(clauses))
	for _, col := range cols {
		vals = append(vals, clauses[col])
	}

	b = builder.Set(b, "Columns", cols).(InsertBuilder)
	b = builder.Set(b, "Values", [][]interface{}{vals}).(InsertBuilder)

	return b
}

// Select set Select clause for insert query
// If Values and Select are used, then Select has higher priority
func (b InsertBuilder) Select(sb selectBuilder) InsertBuilder {
	return builder.Set(b, "Select", &sb).(InsertBuilder)
}

func (b InsertBuilder) statementKeyword(keyword safeString) InsertBuilder {
	return builder.Set(b, "StatementKeyword", keyword).(InsertBuilder)
}
