package squirrel

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
)

type selectData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	Options           []safeString
	Columns           []Sqlizer
	From              Sqlizer
	Joins             []Sqlizer
	WhereParts        []Sqlizer
	GroupBys          []safeString
	HavingParts       []Sqlizer
	OrderByParts      []Sqlizer
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

func (d *selectData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return ExecWith(d.RunWith, d)
}

func (d *selectData) Query() (*sql.Rows, error) {
	if d.RunWith == nil {
		return nil, ErrRunnerNotSet
	}
	return QueryWith(d.RunWith, d)
}

func (d *selectData) QueryRow() RowScanner {
	if d.RunWith == nil {
		return &Row{err: ErrRunnerNotSet}
	}
	queryRower, ok := d.RunWith.(QueryRower)
	if !ok {
		return &Row{err: ErrRunnerNotQueryRunner}
	}
	return QueryRowWith(queryRower, d)
}

func (d *selectData) ToSql() (sqlStr string, args []interface{}, err error) {
	sqlStr, args, err = d.toSqlRaw()
	if err != nil {
		return
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sqlStr)
	return
}

func (d *selectData) toSqlRaw() (sqlStr string, args []interface{}, err error) {
	if len(d.Columns) == 0 {
		err = fmt.Errorf("select statements must have at least one result column")
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

	sql.WriteString("SELECT ")

	if len(d.Options) > 0 {
		for _, val := range d.Options {
			sql.WriteString(string(val))
			sql.WriteString(" ")
		}
	}

	if len(d.Columns) > 0 {
		args, err = appendToSql(d.Columns, sql, ", ", args)
		if err != nil {
			return
		}
	}

	if d.From != nil {
		sql.WriteString(" FROM ")
		args, err = appendToSql([]Sqlizer{d.From}, sql, "", args)
		if err != nil {
			return
		}
	}

	if len(d.Joins) > 0 {
		sql.WriteString(" ")
		args, err = appendToSql(d.Joins, sql, " ", args)
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

	if len(d.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		for idx, val := range d.GroupBys {
			if idx != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(string(val))
		}
	}

	if len(d.HavingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSql(d.HavingParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderByParts) > 0 {
		sql.WriteString(" ORDER BY ")
		args, err = appendToSql(d.OrderByParts, sql, ", ", args)
		if err != nil {
			return
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

	sqlStr = sql.String()
	return
}

// Builder

// selectBuilder builds SQL SELECT statements.
type selectBuilder struct {
	data selectData
}

func SelectBuilder() selectBuilder {
	return selectBuilder{
		data: selectData{
			PlaceholderFormat: Question,
			Prefixes:          make([]Sqlizer, 0),
			Options:           make([]safeString, 0),
			Columns:           make([]Sqlizer, 0),
			Joins:             make([]Sqlizer, 0),
			WhereParts:        make([]Sqlizer, 0),
			GroupBys:          make([]safeString, 0),
			HavingParts:       make([]Sqlizer, 0),
			OrderByParts:      make([]Sqlizer, 0),
			Suffixes:          make([]Sqlizer, 0),
		},
	}
}

// Commenting this for testing direct builder
// type selectBuilder builder.Builder

// func init() {
// 	builder.Register(selectBuilder{}, selectData{})
// }

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b selectBuilder) PlaceholderFormat(f PlaceholderFormat) selectBuilder {
	b.data.PlaceholderFormat = f
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
// For most cases runner will be a database connection.
//
// Internally we use this to mock out the database connection for testing.
func (b selectBuilder) RunWith(runner BaseRunner) selectBuilder {
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
func (b selectBuilder) Exec() (sql.Result, error) {
	return b.data.Exec()
}

// Query builds and Querys the query with the Runner set by RunWith.
func (b selectBuilder) Query() (*sql.Rows, error) {
	return b.data.Query()
}

// QueryRow builds and QueryRows the query with the Runner set by RunWith.
func (b selectBuilder) QueryRow() RowScanner {
	return b.data.QueryRow()
}

// Scan is a shortcut for QueryRow().Scan.
func (b selectBuilder) Scan(dest ...interface{}) error {
	return b.QueryRow().Scan(dest...)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b selectBuilder) ToSql() (string, []interface{}, error) {
	return b.data.ToSql()
}

func (b selectBuilder) toSqlRaw() (string, []interface{}, error) {
	return b.data.toSqlRaw()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b selectBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// Prefix adds an expression to the beginning of the query
func (b selectBuilder) Prefix(sql safeString, args ...interface{}) selectBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b selectBuilder) PrefixExpr(expr Sqlizer) selectBuilder {
	b.data.Prefixes = append(b.data.Prefixes, expr)
	return b
}

// Distinct adds a DISTINCT clause to the query.
func (b selectBuilder) Distinct() selectBuilder {
	return b.Options("DISTINCT")
}

// Options adds select option to the query
func (b selectBuilder) Options(options ...safeString) selectBuilder {
	b.data.Options = append(b.data.Options, options...)
	return b
}

// Columns adds result columns to the query.
func (b selectBuilder) Columns(columns ...safeString) selectBuilder {
	for _, str := range columns {
		b.data.Columns = append(b.data.Columns, str)
	}
	return b
}

// RemoveColumns remove all columns from query.
// Must add a new column with Column or Columns methods, otherwise
// return a error.
func (b selectBuilder) RemoveColumns() selectBuilder {
	b.data.Columns = []Sqlizer{}
	return b
}

// Column adds a result column to the query.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	Column("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3)
func (b selectBuilder) Column(expr Sqlizer) selectBuilder {
	b.data.Columns = append(b.data.Columns, expr)
	return b
}

// From sets the FROM clause of the query.
func (b selectBuilder) From(from safeString) selectBuilder {
	b.data.From = from
	return b
}

// FromSelect sets a subquery into the FROM clause of the query.
func (b selectBuilder) FromSelect(from selectBuilder, alias safeString) selectBuilder {
	// Prevent misnumbered parameters in nested selects (#183).
	from = from.PlaceholderFormat(Question)
	b.data.From = Alias(from, alias)
	return b
}

// JoinClause adds a join clause to the query.
func (b selectBuilder) JoinClause(expr Sqlizer) selectBuilder {
	b.data.Joins = append(b.data.Joins, expr)
	return b
}

// Join adds a JOIN clause to the query.
func (b selectBuilder) Join(join safeString, rest ...interface{}) selectBuilder {
	return b.JoinClause(Expr("JOIN "+join, rest...))
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b selectBuilder) LeftJoin(join safeString, rest ...interface{}) selectBuilder {
	return b.JoinClause(Expr("LEFT JOIN "+join, rest...))
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b selectBuilder) RightJoin(join safeString, rest ...interface{}) selectBuilder {
	return b.JoinClause(Expr("RIGHT JOIN "+join, rest...))
}

// InnerJoin adds a INNER JOIN clause to the query.
func (b selectBuilder) InnerJoin(join safeString, rest ...interface{}) selectBuilder {
	return b.JoinClause(Expr("INNER JOIN "+join, rest...))
}

// CrossJoin adds a CROSS JOIN clause to the query.
func (b selectBuilder) CrossJoin(join safeString, rest ...interface{}) selectBuilder {
	return b.JoinClause(Expr("CROSS JOIN "+join, rest...))
}

// Where adds an expression to the WHERE clause of the query.
//
// Expressions are ANDed together in the generated SQL.
//
// Where accepts several types for its pred argument:
//
// nil OR "" - ignored.
//
// string - SQL expression.
// If the expression has SQL placeholders then a set of arguments must be passed
// as well, one for each placeholder.
//
// map[string]interface{} OR Eq - map of SQL expressions to values. Each key is
// transformed into an expression like "<key> = ?", with the corresponding value
// bound to the placeholder. If the value is nil, the expression will be "<key>
// IS NULL". If the value is an array or slice, the expression will be "<key> IN
// (?,?,...)", with one placeholder for each item in the value. These expressions
// are ANDed together.
//
// Where will panic if pred isn't any of the above types.
func (b selectBuilder) Where(expr Sqlizer) selectBuilder {
	if expr == nil {
		return b
	}
	b.data.WhereParts = append(b.data.WhereParts, expr)
	return b
}

// GroupBy adds GROUP BY expressions to the query.
func (b selectBuilder) GroupBy(groupBys ...safeString) selectBuilder {
	b.data.GroupBys = append(b.data.GroupBys, groupBys...)
	return b
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b selectBuilder) Having(expr Sqlizer) selectBuilder {
	if expr == nil {
		return b
	}
	b.data.HavingParts = append(b.data.HavingParts, expr)
	return b
}

// OrderByClause adds ORDER BY clause to the query.
func (b selectBuilder) OrderByClause(expr Sqlizer) selectBuilder {
	b.data.OrderByParts = append(b.data.OrderByParts, expr)
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b selectBuilder) OrderBy(orderBys ...safeString) selectBuilder {
	for _, orderBy := range orderBys {
		b = b.OrderByClause(orderBy)
	}
	return b
}

// Limit sets a LIMIT clause on the query.
func (b selectBuilder) Limit(limit uint64) selectBuilder {
	b.data.Limit = strconv.FormatUint(limit, 10)
	return b
}

// Limit ALL allows to access all records with limit
func (b selectBuilder) RemoveLimit() selectBuilder {
	b.data.Limit = ""
	return b
}

// Offset sets a OFFSET clause on the query.
func (b selectBuilder) Offset(offset uint64) selectBuilder {
	b.data.Offset = strconv.FormatUint(offset, 10)
	return b
}

// RemoveOffset removes OFFSET clause.
func (b selectBuilder) RemoveOffset() selectBuilder {
	b.data.Offset = ""
	return b
}

// Suffix adds an expression to the end of the query
func (b selectBuilder) Suffix(sql safeString, args ...interface{}) selectBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b selectBuilder) SuffixExpr(expr Sqlizer) selectBuilder {
	b.data.Suffixes = append(b.data.Suffixes, expr)
	return b
}
