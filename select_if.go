package squirrel

import (
	"fmt"

	"github.com/lann/builder"
)

// PlaceholderFormatIf sets PlaceholderFormat (e.g. Question or Dollar) for the
// query if include is true.
func (b selectBuilder) PlaceholderFormatIf(f PlaceholderFormat, include bool) selectBuilder {
	if include {
		return b.PlaceholderFormat(f)
	}
	return b
}

// RunWithIf sets a Runner (like database/sql.DB) to be used with e.g. Exec if include is true.
// For most cases runner will be a database connection.
//
// Internally we use this to mock out the database connection for testing.
func (b selectBuilder) RunWithIf(runner BaseRunner, include bool) selectBuilder {
	if include {
		return b.RunWith(runner)
	}
	return b
}

// PrefixIf adds an expression to the beginning of the query if include is true.
func (b selectBuilder) PrefixIf(sql safeString, include bool, args ...interface{}) selectBuilder {
	if include {
		return b.PrefixExpr(Expr(sql, args...))
	}
	return b
}

// PrefixExprIf adds an expression to the very beginning of the query if include is true.
func (b selectBuilder) PrefixExprIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return builder.Append(b, "Prefixes", expr).(selectBuilder)
	}
	return b
}

// DistinctIf adds a DISTINCT clause to the query if include is true.
func (b selectBuilder) DistinctIf(include bool) selectBuilder {
	if include {
		return b.Options("DISTINCT")
	}
	return b
}

// OptionsIf adds select option to the query for each Include that is true.
func (b selectBuilder) OptionsIf(options ...valIf[safeString]) selectBuilder {
	opts := make([]safeString, 0)
	for _, v := range options {
		if v.Include {
			opts = append(opts, v.Value)
		}
	}
	return builder.Extend(b, "Options", opts).(selectBuilder)
}

// OptionIf adds a single select option to the query if include is true.
func (b selectBuilder) OptionIf(option safeString, include bool) selectBuilder {
	if include {
		return builder.Append(b, "Options", option).(selectBuilder)
	}
	return b
}

// ColumnsIf adds result columns to the query for each Include that is true.
func (b selectBuilder) ColumnsIf(columns ...valIf[safeString]) selectBuilder {
	parts := make([]Sqlizer, 0, len(columns))
	for _, v := range columns {
		if v.Include {
			parts = append(parts, v.Value)
		}
	}
	return builder.Extend(b, "Columns", parts).(selectBuilder)
}

// ColumnIf adds a result column to the query if include is true.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	ColumnIf(Expr("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), true) == "IF(col IN (1, 2, 3), 1, 0) as col"
//	ColumnIf(Expr("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), false) == ""
func (b selectBuilder) ColumnIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return builder.Append(b, "Columns", expr).(selectBuilder)
	}
	return b
}

// JoinClauseIf adds a join clause to the query if include is true.
func (b selectBuilder) JoinClauseIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return builder.Append(b, "Joins", expr).(selectBuilder)
	}
	return b
}

// JoinIf adds a JOIN clause to the query if include is true.
func (b selectBuilder) JoinIf(join safeString, include bool, rest ...interface{}) selectBuilder {
	if include {
		return b.JoinClause(Expr("JOIN "+join, rest...))
	}
	return b
}

// LeftJoinIf adds a LEFT JOIN clause to the query if include is true.
func (b selectBuilder) LeftJoinIf(join safeString, include bool, rest ...interface{}) selectBuilder {
	if include {
		return b.JoinClause(Expr("LEFT JOIN "+join, rest...))
	}
	return b
}

// RightJoinIf adds a RIGHT JOIN clause to the query if include is true.
func (b selectBuilder) RightJoinIf(join safeString, include bool, rest ...interface{}) selectBuilder {
	if include {
		return b.JoinClause(Expr("RIGHT JOIN "+join, rest...))
	}
	return b
}

// InnerJoinIf adds a INNER JOIN clause to the query if include is true.
func (b selectBuilder) InnerJoinIf(join safeString, include bool, rest ...interface{}) selectBuilder {
	if include {
		return b.JoinClause(Expr("INNER JOIN "+join, rest...))
	}
	return b
}

// CrossJoinIf adds a CROSS JOIN clause to the query if include is true.
func (b selectBuilder) CrossJoinIf(join safeString, include bool, rest ...interface{}) selectBuilder {
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
func (b selectBuilder) WhereIf(expr Sqlizer, include bool) selectBuilder {
	if include && expr != nil {
		return builder.Append(b, "WhereParts", expr).(selectBuilder)
	}
	return b
}

// GroupBysIf adds GROUP BY expressions to the query for each Include that is true.
func (b selectBuilder) GroupBysIf(groupBys ...valIf[safeString]) selectBuilder {
	grps := make([]safeString, 0)
	for _, v := range groupBys {
		if v.Include {
			grps = append(grps, v.Value)
		}
	}
	return builder.Extend(b, "GroupBys", grps).(selectBuilder)
}

// GroupByIf adds a single GROUP BY expression to the query if include is true.
func (b selectBuilder) GroupByIf(groupBy safeString, include bool) selectBuilder {
	if include {
		return builder.Append(b, "GroupBys", groupBy).(selectBuilder)
	}
	return b
}

// HavingIf adds an expression to the HAVING clause of the query if include is true.
//
// See Where.
func (b selectBuilder) HavingIf(expr Sqlizer, include bool) selectBuilder {
	if include && expr != nil {
		return builder.Append(b, "HavingParts", expr).(selectBuilder)
	}
	return b
}

// OrderByClauseIf adds ORDER BY clause to the query if include is true.
func (b selectBuilder) OrderByClauseIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return builder.Append(b, "OrderByParts", expr).(selectBuilder)
	}
	return b
}

// OrderBysIf adds ORDER BY expressions to the query for each Include that is true.
func (b selectBuilder) OrderBysIf(orderBys ...valIf[safeString]) selectBuilder {
	for _, orderBy := range orderBys {
		if orderBy.Include {
			b = b.OrderByClause(orderBy.Value)
		}
	}

	return b
}

// OrderByIf adds a single ORDER BY expression to the query if include is true.
func (b selectBuilder) OrderByIf(orderBy safeString, include bool) selectBuilder {
	if include {
		b = b.OrderByClause(orderBy)
	}
	return b
}

// LimitIf sets a LIMIT clause on the query if include is true.
func (b selectBuilder) LimitIf(limit uint64, include bool) selectBuilder {
	if include {
		return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(selectBuilder)
	}
	return b
}

// OffsetIf sets a OFFSET clause on the query if include is true.
func (b selectBuilder) OffsetIf(offset uint64, include bool) selectBuilder {
	if include {
		return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(selectBuilder)
	}
	return b
}

// SuffixIf adds an expression to the end of the query if include is true.
func (b selectBuilder) SuffixIf(sql safeString, include bool, args ...interface{}) selectBuilder {
	if include {
		return b.SuffixExpr(Expr(sql, args...))
	}
	return b
}

// SuffixExprIf adds an expression to the end of the query if include is true.
func (b selectBuilder) SuffixExprIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return b.SuffixExpr(expr)
	}
	return b
}
