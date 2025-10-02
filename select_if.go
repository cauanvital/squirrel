package squirrel2

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
		return b.PrefixExpr(expr)
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
	for _, v := range options {
		b = b.OptionIf(v.Value, v.Include)
	}
	return b
}

// OptionIf adds a single select option to the query if include is true.
func (b selectBuilder) OptionIf(option safeString, include bool) selectBuilder {
	if include {
		return b.Options(option)
	}
	return b
}

// ColumnsIf adds result columns to the query for each Include that is true.
func (b selectBuilder) ColumnsIf(columns ...valIf[safeString]) selectBuilder {
	for _, v := range columns {
		b = b.ColumnIf(v.Value, v.Include)
	}
	return b
}

// ColumnIf adds a result column to the query if include is true.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	ColumnIf(Expr("IF(col IN ("+squirrel2.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), true) == "IF(col IN (1, 2, 3), 1, 0) as col"
//	ColumnIf(Expr("IF(col IN ("+squirrel2.Placeholders(3)+"), 1, 0) as col", 1, 2, 3), false) == ""
func (b selectBuilder) ColumnIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return b.Column(expr)
	}
	return b
}

// JoinClauseIf adds a join clause to the query if include is true.
func (b selectBuilder) JoinClauseIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return b.JoinClause(expr)
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
	if include {
		return b.Where(expr)
	}
	return b
}

// GroupBysIf adds GROUP BY expressions to the query for each Include that is true.
func (b selectBuilder) GroupBysIf(groupBys ...valIf[safeString]) selectBuilder {
	for _, v := range groupBys {
		b = b.GroupByIf(v.Value, v.Include)
	}
	return b
}

// GroupByIf adds a single GROUP BY expression to the query if include is true.
func (b selectBuilder) GroupByIf(groupBy safeString, include bool) selectBuilder {
	if include {
		return b.GroupBy(groupBy)
	}
	return b
}

// HavingIf adds an expression to the HAVING clause of the query if include is true.
//
// See Where.
func (b selectBuilder) HavingIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return b.Having(expr)
	}
	return b
}

// OrderByClauseIf adds ORDER BY clause to the query if include is true.
func (b selectBuilder) OrderByClauseIf(expr Sqlizer, include bool) selectBuilder {
	if include {
		return b.OrderByClause(expr)
	}
	return b
}

// OrderBysIf adds ORDER BY expressions to the query for each Include that is true.
func (b selectBuilder) OrderByIf(orderBys ...valIf[safeString]) selectBuilder {
	for _, v := range orderBys {
		b = b.OrderByClauseIf(v.Value, v.Include)
	}
	return b
}

// LimitIf sets a LIMIT clause on the query if include is true.
func (b selectBuilder) LimitIf(limit uint64, include bool) selectBuilder {
	if include {
		return b.Limit(limit)
	}
	return b
}

// OffsetIf sets a OFFSET clause on the query if include is true.
func (b selectBuilder) OffsetIf(offset uint64, include bool) selectBuilder {
	if include {
		return b.Offset(offset)
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
