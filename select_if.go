package squirrel

import (
	"fmt"

	"github.com/lann/builder"
)

// PlaceholderFormatIf sets PlaceholderFormat (e.g. Question or Dollar) for the
// query if include is true.
func (b SelectBuilder) PlaceholderFormatIf(f PlaceholderFormat, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "PlaceholderFormat", f).(SelectBuilder)
	}
	return b
}

// RunWithIf sets a Runner (like database/sql.DB) to be used with e.g. Exec if include is true.
// For most cases runner will be a database connection.
//
// Internally we use this to mock out the database connection for testing.
func (b SelectBuilder) RunWithIf(runner BaseRunner, include bool) SelectBuilder {
	if include {
		return setRunWith(b, runner).(SelectBuilder)
	}
	return b
}

// PrefixIf adds an expression to the beginning of the query if include is true.
func (b SelectBuilder) PrefixIf(sql safeString, include bool, args ...interface{}) SelectBuilder {
	if include {
		return b.PrefixExpr(Expr(sql, args...))
	}
	return b
}

// PrefixExprIf adds an expression to the very beginning of the query if include is true.
func (b SelectBuilder) PrefixExprIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Prefixes", expr).(SelectBuilder)
	}
	return b
}

// DistinctIf adds a DISTINCT clause to the query if include is true.
func (b SelectBuilder) DistinctIf(include bool) SelectBuilder {
	if include {
		return b.Options("DISTINCT")
	}
	return b
}

// OptionsIf adds select option to the query for each Include that is true.
func (b SelectBuilder) OptionsIf(options ...ValIf[safeString]) SelectBuilder {
	opts := make([]safeString, 0)
	for _, v := range options {
		if v.Include {
			opts = append(opts, v.Value)
		}
	}
	return builder.Extend(b, "Options", opts).(SelectBuilder)
}

// OptionIf adds a single select option to the query if include is true.
func (b SelectBuilder) OptionIf(option safeString, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Options", option).(SelectBuilder)
	}
	return b
}

// ColumnsIf adds result columns to the query for each Include that is true.
func (b SelectBuilder) ColumnsIf(columns ...ValIf[safeString]) SelectBuilder {
	parts := make([]Sqlizer, 0, len(columns))
	for _, v := range columns {
		if v.Include {
			parts = append(parts, v.Value)
		}
	}
	return builder.Extend(b, "Columns", parts).(SelectBuilder)
}

// ColumnIf adds a result column to the query if include is true.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	ColumnIf(Expr("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), true) == "IF(col IN (1, 2, 3), 1, 0) as col"
//	ColumnIf(Expr("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), false) == ""
func (b SelectBuilder) ColumnIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Columns", expr).(SelectBuilder)
	}
	return b
}

// JoinClauseIf adds a join clause to the query if include is true.
func (b SelectBuilder) JoinClauseIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Joins", expr).(SelectBuilder)
	}
	return b
}

// JoinIf adds a JOIN clause to the query if include is true.
func (b SelectBuilder) JoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("JOIN "+join, rest...))
	}
	return b
}

// LeftJoinIf adds a LEFT JOIN clause to the query if include is true.
func (b SelectBuilder) LeftJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("LEFT JOIN "+join, rest...))
	}
	return b
}

// RightJoinIf adds a RIGHT JOIN clause to the query if include is true.
func (b SelectBuilder) RightJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("RIGHT JOIN "+join, rest...))
	}
	return b
}

// InnerJoinIf adds a INNER JOIN clause to the query if include is true.
func (b SelectBuilder) InnerJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("INNER JOIN "+join, rest...))
	}
	return b
}

// CrossJoinIf adds a CROSS JOIN clause to the query if include is true.
func (b SelectBuilder) CrossJoinIf(join safeString, include bool, rest ...interface{}) SelectBuilder {
	if include {
		return b.JoinClause(Expr("CROSS JOIN "+join, rest...))
	}
	return b
}

// WhereIf adds an expression to the WHERE clause of the query if include is true.
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
func (b SelectBuilder) WhereIf(expr Sqlizer, include bool) SelectBuilder {
	if include && expr != nil {
		return builder.Append(b, "WhereParts", expr).(SelectBuilder)
	}
	return b
}

// GroupBysIf adds GROUP BY expressions to the query for each Include that is true.
func (b SelectBuilder) GroupBysIf(groupBys ...ValIf[safeString]) SelectBuilder {
	grps := make([]safeString, 0)
	for _, v := range groupBys {
		if v.Include {
			grps = append(grps, v.Value)
		}
	}
	return builder.Extend(b, "GroupBys", grps).(SelectBuilder)
}

// GroupByIf adds a single GROUP BY expression to the query if include is true.
func (b SelectBuilder) GroupByIf(groupBy safeString, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "GroupBys", groupBy).(SelectBuilder)
	}
	return b
}

// HavingIf adds an expression to the HAVING clause of the query if include is true.
//
// See Where.
func (b SelectBuilder) HavingIf(expr Sqlizer, include bool) SelectBuilder {
	if include && expr != nil {
		return builder.Append(b, "HavingParts", expr).(SelectBuilder)
	}
	return b
}

// OrderByClauseIf adds ORDER BY clause to the query if include is true.
func (b SelectBuilder) OrderByClauseIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "OrderByParts", expr).(SelectBuilder)
	}
	return b
}

// OrderBysIf adds ORDER BY expressions to the query for each Include that is true.
func (b SelectBuilder) OrderBysIf(orderBys ...ValIf[safeString]) SelectBuilder {
	for _, orderBy := range orderBys {
		if orderBy.Include {
			b = b.OrderByClause(orderBy.Value)
		}
	}

	return b
}

// OrderByIf adds a single ORDER BY expression to the query if include is true.
func (b SelectBuilder) OrderByIf(orderBy safeString, include bool) SelectBuilder {
	if include {
		b = b.OrderByClause(orderBy)
	}
	return b
}

// LimitIf sets a LIMIT clause on the query if include is true.
func (b SelectBuilder) LimitIf(limit uint64, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(SelectBuilder)
	}
	return b
}

// OffsetIf sets a OFFSET clause on the query if include is true.
func (b SelectBuilder) OffsetIf(offset uint64, include bool) SelectBuilder {
	if include {
		return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(SelectBuilder)
	}
	return b
}

// SuffixIf adds an expression to the end of the query if include is true.
func (b SelectBuilder) SuffixIf(sql safeString, include bool, args ...interface{}) SelectBuilder {
	if include {
		return b.SuffixExpr(Expr(sql, args...))
	}
	return b
}

// SuffixExprIf adds an expression to the end of the query if include is true.
func (b SelectBuilder) SuffixExprIf(expr Sqlizer, include bool) SelectBuilder {
	if include {
		return builder.Append(b, "Suffixes", expr).(SelectBuilder)
	}
	return b
}
