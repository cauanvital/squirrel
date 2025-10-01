package squirrel

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type updateData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	Table             safeString
	SetClauses        []setClause
	From              Sqlizer
	WhereParts        []Sqlizer
	OrderBys          []safeString
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

type setClause struct {
	column safeString
	value  interface{}
}

func (d *updateData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *updateData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}

func (d *updateData) QueryRow() RowScanner {
	if d.RunWith == nil {
		return &Row{err: ErrRunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRower)
	if !ok {
		return &Row{err: ErrRunnerNotQueryRunner}
	}
	return QueryRowWith(queryRower, d)
}

func (d *updateData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.Table) == 0 {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(d.SetClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set clause")
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

	sql.WriteString("UPDATE ")
	sql.WriteString(string(d.Table))

	sql.WriteString(" SET ")
	setSqls := make([]string, len(d.SetClauses))
	for i, setClause := range d.SetClauses {
		var valSql string
		if vs, ok := setClause.value.(Sqlizer); ok {
			vsql, vargs, err := vs.ToSql()
			if err != nil {
				return "", nil, err
			}
			if _, ok := vs.(selectBuilder); ok {
				valSql = fmt.Sprintf("(%s)", vsql)
			} else {
				valSql = vsql
			}
			args = append(args, vargs...)
		} else {
			valSql = "?"
			args = append(args, setClause.value)
		}
		setSqls[i] = fmt.Sprintf("%s = %s", setClause.column, valSql)
	}
	sql.WriteString(strings.Join(setSqls, ", "))

	if d.From != nil {
		sql.WriteString(" FROM ")
		args, err = appendToSql([]Sqlizer{d.From}, sql, "", args)
		if err != nil {
			return
		}
	}

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

// UpdateBuilder builds SQL UPDATE statements.
type updateBuilder struct {
	data updateData
}

func UpdateBuilder(b statementBuilderType) updateBuilder {
	return updateBuilder{
		data: updateData{
			PlaceholderFormat: b.placeholderFormat,
			RunWith:           b.runWith,
			WhereParts:        b.whereParts,
			Prefixes:          make([]Sqlizer, 0),
			SetClauses:        make([]setClause, 0),
			OrderBys:          make([]safeString, 0),
			Suffixes:          make([]Sqlizer, 0),
		},
	}
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b updateBuilder) PlaceholderFormat(f PlaceholderFormat) updateBuilder {
	b.data.PlaceholderFormat = f
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b updateBuilder) RunWith(runner BaseRunner) updateBuilder {
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
func (b updateBuilder) Exec() (sql.Result, error) {
	return b.data.Exec()
}

func (b updateBuilder) Query() (*sql.Rows, error) {
	return b.data.Query()
}

func (b updateBuilder) QueryRow() RowScanner {
	return b.data.QueryRow()
}

func (b updateBuilder) Scan(dest ...interface{}) error {
	return b.QueryRow().Scan(dest...)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b updateBuilder) ToSql() (string, []interface{}, error) {
	return b.data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b updateBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query
func (b updateBuilder) Prefix(sql safeString, args ...interface{}) updateBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b updateBuilder) PrefixExpr(expr Sqlizer) updateBuilder {
	b.data.Prefixes = append(b.data.Prefixes, expr)
	return b
}

// Table sets the table to be updated.
func (b updateBuilder) Table(table safeString) updateBuilder {
	b.data.Table = table
	return b
}

// Set adds SET clauses to the query.
func (b updateBuilder) Set(column safeString, value interface{}) updateBuilder {
	b.data.SetClauses = append(b.data.SetClauses, setClause{column: column, value: value})
	return b
}

// SetMap is a convenience method which calls .Set for each key/value pair in clauses.
func (b updateBuilder) SetMap(clauses map[safeString]interface{}) updateBuilder {
	keys := make([]safeString, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}
	sort.Slice(keys, func(idx, jdx int) bool {
		return keys[idx] < keys[jdx]
	})
	for _, key := range keys {
		val := clauses[key]
		b = b.Set(key, val)
	}
	return b
}

// From adds FROM clause to the query
// FROM is valid construct in postgresql only.
func (b updateBuilder) From(from safeString) updateBuilder {
	b.data.From = from
	return b
}

// FromSelect sets a subquery into the FROM clause of the query.
func (b updateBuilder) FromSelect(from selectBuilder, alias safeString) updateBuilder {
	b.data.From = Alias(from.PlaceholderFormat(Question), alias)
	return b
}

// Where adds WHERE expressions to the query.
//
// See selectBuilder.Where for more information.
func (b updateBuilder) Where(expr Sqlizer) updateBuilder {
	b.data.WhereParts = append(b.data.WhereParts, expr)
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b updateBuilder) OrderBy(orderBys ...safeString) updateBuilder {
	b.data.OrderBys = append(b.data.OrderBys, orderBys...)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b updateBuilder) Limit(limit uint64) updateBuilder {
	b.data.Limit = strconv.FormatUint(limit, 10)
	return b
}

// Offset sets a OFFSET clause on the query.
func (b updateBuilder) Offset(offset uint64) updateBuilder {
	b.data.Offset = strconv.FormatUint(offset, 10)
	return b
}

// Suffix adds an expression to the end of the query
func (b updateBuilder) Suffix(sql safeString, args ...interface{}) updateBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b updateBuilder) SuffixExpr(expr Sqlizer) updateBuilder {
	b.data.Suffixes = append(b.data.Suffixes, expr)
	return b
}
