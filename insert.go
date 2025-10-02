package squirrel2

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type insertData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	StatementKeyword  safeString
	Options           []safeString
	Into              safeString
	Columns           []safeString
	Values            [][]interface{}
	Suffixes          []Sqlizer
	Select            *selectBuilder
}

func (d *insertData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *insertData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}

func (d *insertData) QueryRow() RowScanner {
	if d.RunWith == nil {
		return &Row{err: ErrRunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRower)
	if !ok {
		return &Row{err: ErrRunnerNotQueryRunner}
	}
	return QueryRowWith(queryRower, d)
}

func (d *insertData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.Into) == 0 {
		err = errors.New("insert statements must specify a table")
		return
	}
	if len(d.Values) == 0 && d.Select == nil {
		err = errors.New("insert statements must have at least one set of values or select clause")
		return
	}

	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSql(d.Prefixes, sql, " ", args)
		if err != nil {
			return
		}

		sql.WriteString(" ")
	}

	if d.StatementKeyword == "" {
		sql.WriteString("INSERT ")
	} else {
		sql.WriteString(string(d.StatementKeyword))
		sql.WriteString(" ")
	}

	if len(d.Options) > 0 {
		for _, val := range d.Options {
			sql.WriteString(string(val))
			sql.WriteString(" ")
		}
	}

	sql.WriteString("INTO ")
	sql.WriteString(string(d.Into))
	sql.WriteString(" ")

	if len(d.Columns) > 0 {
		sql.WriteString("(")
		for idx, val := range d.Columns {
			if idx != 0 {
				sql.WriteString(",")
			}
			sql.WriteString(string(val))
		}
		sql.WriteString(") ")
	}

	if d.Select != nil {
		args, err = d.appendSelectToSQL(sql, args)
	} else {
		args, err = d.appendValuesToSQL(sql, args)
	}
	if err != nil {
		return
	}

	if len(d.Suffixes) > 0 {
		sql.WriteString(" ")
		args, err = appendToSql(d.Suffixes, sql, " ", args)
		if err != nil {
			return
		}
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sql.String())
	return
}

func (d *insertData) appendValuesToSQL(w io.Writer, args []interface{}) ([]interface{}, error) {
	if len(d.Values) == 0 {
		return args, errors.New("values for insert statements are not set")
	}

	io.WriteString(w, "VALUES ")

	valuesStrings := make([]string, len(d.Values))
	for r, row := range d.Values {
		valueStrings := make([]string, len(row))
		for v, val := range row {
			if vs, ok := val.(Sqlizer); ok {
				vsql, vargs, err := vs.ToSql()
				if err != nil {
					return nil, err
				}
				valueStrings[v] = vsql
				args = append(args, vargs...)
			} else {
				valueStrings[v] = "?"
				args = append(args, val)
			}
		}
		valuesStrings[r] = fmt.Sprintf("(%s)", strings.Join(valueStrings, ","))
	}

	io.WriteString(w, strings.Join(valuesStrings, ","))

	return args, nil
}

func (d *insertData) appendSelectToSQL(w io.Writer, args []interface{}) ([]interface{}, error) {
	if d.Select == nil {
		return args, errors.New("select clause for insert statements are not set")
	}

	selectClause, sArgs, err := d.Select.ToSql()
	if err != nil {
		return args, err
	}

	io.WriteString(w, selectClause)
	args = append(args, sArgs...)

	return args, nil
}

// Builder

// InsertBuilder builds SQL INSERT statements.
type insertBuilder struct {
	data insertData
}

func InsertBuilder(b statementBuilderType) insertBuilder {
	return insertBuilder{
		data: insertData{
			PlaceholderFormat: b.placeholderFormat,
			RunWith:           b.runWith,
			Prefixes:          make([]Sqlizer, 0),
			Options:           make([]safeString, 0),
			Columns:           make([]safeString, 0),
			Values:            make([][]interface{}, 0),
			Suffixes:          make([]Sqlizer, 0),
		},
	}
}

// Commenting this for testing direct builder
// type InsertBuilder builder.Builder

// func init() {
// 	builder.Register(InsertBuilder{}, insertData{})
// }

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b insertBuilder) PlaceholderFormat(f PlaceholderFormat) insertBuilder {
	b.data.PlaceholderFormat = f
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b insertBuilder) RunWith(runner BaseRunner) insertBuilder {
	switch r := runner.(type) {
	case StdSqlCtx:
		runner = WrapStdSqlCtx(r)
	case StdSql:
		runner = WrapStdSql(r)
	}
	b.data.RunWith = runner
	return b
}

// Exec builds and Execs the query with the Runner set by RunWith.
func (b insertBuilder) Exec() (sql.Result, error) {
	return b.data.Exec()
}

// Query builds and Querys the query with the Runner set by RunWith.
func (b insertBuilder) Query() (*sql.Rows, error) {
	return b.data.Query()
}

// QueryRow builds and QueryRows the query with the Runner set by RunWith.
func (b insertBuilder) QueryRow() RowScanner {
	return b.data.QueryRow()
}

// Scan is a shortcut for QueryRow().Scan.
func (b insertBuilder) Scan(dest ...interface{}) error {
	return b.QueryRow().Scan(dest...)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b insertBuilder) ToSql() (string, []interface{}, error) {
	return b.data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b insertBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query
func (b insertBuilder) Prefix(sql safeString, args ...interface{}) insertBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b insertBuilder) PrefixExpr(expr Sqlizer) insertBuilder {
	b.data.Prefixes = append(b.data.Prefixes, expr)
	return b
}

// Options adds keyword options before the INTO clause of the query.
func (b insertBuilder) Options(options ...safeString) insertBuilder {
	b.data.Options = append(b.data.Options, options...)
	return b
}

// Into sets the INTO clause of the query.
func (b insertBuilder) Into(into safeString) insertBuilder {
	b.data.Into = into
	return b
}

// Columns adds insert columns to the query.
func (b insertBuilder) Columns(columns ...safeString) insertBuilder {
	b.data.Columns = append(b.data.Columns, columns...)
	return b
}

// Values adds a single row's values to the query.
func (b insertBuilder) Values(values ...interface{}) insertBuilder {
	b.data.Values = append(b.data.Values, values)
	return b
}

// Suffix adds an expression to the end of the query
func (b insertBuilder) Suffix(sql safeString, args ...interface{}) insertBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b insertBuilder) SuffixExpr(expr Sqlizer) insertBuilder {
	b.data.Suffixes = append(b.data.Suffixes, expr)
	return b
}

// SetMap set columns and values for insert builder from a map of column name and value
// note that it will reset all previous columns and values was set if any
func (b insertBuilder) SetMap(clauses map[safeString]interface{}) insertBuilder {
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

	b.data.Columns = cols
	b.data.Values = [][]interface{}{vals}

	return b
}

// Select set Select clause for insert query
// If Values and Select are used, then Select has higher priority
func (b insertBuilder) Select(sb selectBuilder) insertBuilder {
	b.data.Select = &sb
	return b
}

func (b insertBuilder) statementKeyword(keyword safeString) insertBuilder {
	b.data.StatementKeyword = keyword
	return b
}
