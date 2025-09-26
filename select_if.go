package squirrel

import (
	"fmt"

	"github.com/lann/builder"
)

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b SelectBuilder) PlaceholderFormatIf(f PlaceholderFormat, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "PlaceholderFormat", f).(SelectBuilder)
	}
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
// For most cases runner will be a database connection.
//
// Internally we use this to mock out the database connection for testing.
func (b SelectBuilder) RunWithIf(runner BaseRunner, include bool) SelectBuilder {
	if include {
		return setRunWith(b, runner).(SelectBuilder)
	}
	return b
}

// Prefix adds an expression to the beginning of the query
func (b SelectBuilder) PrefixIf(sql safeString, include bool, args ...interface{}) SelectBuilder {
	if include {
		return b.PrefixExpr(Expr(sql, args...))
	}
	return b
}

// PrefixExpr adds an expression to the very beginning of the query
func (b SelectBuilder) PrefixExprIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Prefixes", expr).(SelectBuilder)
	}
	return b
}

// Distinct adds a DISTINCT clause to the query.
func (b SelectBuilder) DistinctIf(include bool) SelectBuilder {
	if include {
		return b.Options("DISTINCT")
	}
	return b
}

// Options adds select option to the query
func (b SelectBuilder) OptionsIf(options ...ValIf[safeString]) SelectBuilder {
	opts := make([]safeString, 0)
	for _, v := range options {
		if v.Include {
			opts = append(opts, v.Value)
		}
	}
	return builder.Extend(b, "Options", opts).(SelectBuilder)
}

// Columns adds result columns to the query.
func (b SelectBuilder) ColumnsIf(columns ...ValIf[safeString]) SelectBuilder {
	parts := make([]Sqlizer, 0, len(columns))
	for _, v := range columns {
		if v.Include {
			parts = append(parts, v.Value)
		}
	}
	return builder.Extend(b, "Columns", parts).(SelectBuilder)
}

// Column adds a result column to the query.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	Column("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3)
func (b SelectBuilder) ColumnIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Columns", expr).(SelectBuilder)
	}
	return b
}

// Join adds a JOIN clause to the query.
func (b SelectBuilder) JoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("JOIN "+join, rest...))
	}
	return b
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b SelectBuilder) LeftJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("LEFT JOIN "+join, rest...))
	}
	return b
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b SelectBuilder) RightJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("RIGHT JOIN "+join, rest...))
	}
	return b
}

// InnerJoin adds a INNER JOIN clause to the query.
func (b SelectBuilder) InnerJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("INNER JOIN "+join, rest...))
	}
	return b
}

// CrossJoin adds a CROSS JOIN clause to the query.
func (b SelectBuilder) CrossJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("CROSS JOIN "+join, rest...))
	}
	return b
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
func (b SelectBuilder) Where(expr Sqlizer) SelectBuilder {
	if expr == nil {
		return b
	}
	return builder.Append(b, "WhereParts", expr).(SelectBuilder)
}

// GroupBy adds GROUP BY expressions to the query.
func (b SelectBuilder) GroupById(groupBys ...ValIf[safeString]) SelectBuilder {
	grps := make([]safeString, 0)
	for _, v := range groupBys {
		if v.Include {
			grps = append(grps, v.Value)
		}
	}
	return builder.Extend(b, "GroupBys", grps).(SelectBuilder)
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b SelectBuilder) Having(expr Sqlizer) SelectBuilder {
	if expr == nil {
		return b
	}
	return builder.Append(b, "HavingParts", expr).(SelectBuilder)
}

// OrderByClause adds ORDER BY clause to the query.
func (b SelectBuilder) OrderByClause(expr Sqlizer) SelectBuilder {
	return builder.Append(b, "OrderByParts", expr).(SelectBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b SelectBuilder) OrderByIf(orderBys ...ValIf[safeString]) SelectBuilder {
	for _, orderBy := range orderBys {
		if orderBy.Include {
			b = b.OrderByClause(orderBy.Value)
		}
	}

	return b
}

// Limit sets a LIMIT clause on the query.
func (b SelectBuilder) LimitIf(limit uint64, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(SelectBuilder)
	}
	return b
}

// Offset sets a OFFSET clause on the query.
func (b SelectBuilder) OffsetIf(offset uint64, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(SelectBuilder)
	}
	return b
}

// Suffix adds an expression to the end of the query
func (b SelectBuilder) SuffixIf(sql safeString, include bool, args ...interface{}) SelectBuilder {
	if include {
		return b.SuffixExpr(Expr(sql, args...))
	}
	return b
}

// SuffixExpr adds an expression to the end of the query
func (b SelectBuilder) SuffixExpr(expr Sqlizer) SelectBuilder {
	return builder.Append(b, "Suffixes", expr).(SelectBuilder)
}
