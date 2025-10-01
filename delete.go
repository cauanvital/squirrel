package squirrel

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
)

type deleteData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	From              safeString
	WhereParts        []Sqlizer
	OrderBys          []safeString
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

func (d *deleteData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *deleteData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.From) == 0 {
		err = fmt.Errorf("delete statements must specify a From table")
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

	sql.WriteString("DELETE FROM ")
	sql.WriteString(string(d.From))

	if len(d.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(d.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		for idx, val := range d.OrderBys {
			if idx != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(string(val))
		}
	}

	if len(d.Limit) > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(d.Limit)
	}

	if len(d.Offset) > 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(d.Offset)
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

// Builder

// DeleteBuilder builds SQL DELETE statements.
type deleteBuilder struct {
	data deleteData
}

func DeleteBuilder(b statementBuilderType) deleteBuilder {
	return deleteBuilder{
		data: deleteData{
			PlaceholderFormat: b.placeholderFormat,
			RunWith:           b.runWith,
			WhereParts:        b.whereParts,
			Prefixes:          make([]Sqlizer, 0),
			OrderBys:          make([]safeString, 0),
			Suffixes:          make([]Sqlizer, 0),
		},
	}
}

// Commenting this for testing direct builder
// type DeleteBuilder builder.Builder

// func init() {
// 	builder.Register(DeleteBuilder{}, deleteData{})
// }

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b deleteBuilder) PlaceholderFormat(f PlaceholderFormat) deleteBuilder {
	b.data.PlaceholderFormat = f
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b deleteBuilder) RunWith(runner BaseRunner) deleteBuilder {
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
func (b deleteBuilder) Exec() (sql.Result, error) {
	return b.data.Exec()
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b deleteBuilder) ToSql() (string, []interface{}, error) {
	return b.data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b deleteBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query
func (b deleteBuilder) Prefix(sql safeString, args ...interface{}) deleteBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b deleteBuilder) PrefixExpr(expr Sqlizer) deleteBuilder {
	b.data.Prefixes = append(b.data.Prefixes, expr)
	return b
}

// From sets the table to be deleted from.
func (b deleteBuilder) From(from safeString) deleteBuilder {
	b.data.From = from
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b deleteBuilder) Where(expr Sqlizer) deleteBuilder {
	b.data.WhereParts = append(b.data.WhereParts, expr)
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b deleteBuilder) OrderBy(orderBys ...safeString) deleteBuilder {
	b.data.OrderBys = append(b.data.OrderBys, orderBys...)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b deleteBuilder) Limit(limit uint64) deleteBuilder {
	b.data.Limit = strconv.FormatUint(limit, 10)
	return b
}

// Offset sets a OFFSET clause on the query.
func (b deleteBuilder) Offset(offset uint64) deleteBuilder {
	b.data.Offset = strconv.FormatUint(offset, 10)
	return b
}

// Suffix adds an expression to the end of the query
func (b deleteBuilder) Suffix(sql safeString, args ...interface{}) deleteBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b deleteBuilder) SuffixExpr(expr Sqlizer) deleteBuilder {
	b.data.Suffixes = append(b.data.Suffixes, expr)
	return b
}

func (b deleteBuilder) Query() (*sql.Rows, error) {
	return b.data.Query()
}

func (d *deleteData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}
