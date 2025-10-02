package squirrel2

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b deleteBuilder) PlaceholderFormatIf(f PlaceholderFormat, include bool) deleteBuilder {
	if include {
		b.PlaceholderFormat(f)
	}
	return b
}

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b deleteBuilder) RunWithIf(runner BaseRunner, include bool) deleteBuilder {
	if include {
		return b.RunWith(runner)
	}
	return b
}

// Prefix adds an expression to the beginning of the query
func (b deleteBuilder) PrefixIf(sql safeString, include bool, args ...interface{}) deleteBuilder {
	if include {
		return b.Prefix(sql, args...)
	}
	return b
}

// PrefixExpr adds an expression to the very beginning of the query
func (b deleteBuilder) PrefixExprIf(expr Sqlizer, include bool) deleteBuilder {
	if include {
		return b.PrefixExpr(expr)
	}
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b deleteBuilder) WhereIf(expr Sqlizer, include bool) deleteBuilder {
	if include {
		return b.Where(expr)
	}
	return b
}

// OrderBysIf adds ORDER BY expressions to the query for each Include that is true.
func (b deleteBuilder) OrderBysIf(orderBys ...valIf[safeString]) deleteBuilder {
	for _, orderBy := range orderBys {
		b = b.OrderByIf(orderBy.Value, orderBy.Include)
	}
	return b
}

// OrderByIf adds a single ORDER BY expression to the query if include is true.
func (b deleteBuilder) OrderByIf(orderBy safeString, include bool) deleteBuilder {
	if include {
		return b.OrderBy(orderBy)
	}
	return b
}

// Limit sets a LIMIT clause on the query.
func (b deleteBuilder) LimitIf(limit uint64, include bool) deleteBuilder {
	if include {
		return b.Limit(limit)
	}
	return b
}

// Offset sets a OFFSET clause on the query.
func (b deleteBuilder) OffsetIf(offset uint64, include bool) deleteBuilder {
	if include {
		return b.Offset(offset)
	}
	return b
}

// Suffix adds an expression to the end of the query
func (b deleteBuilder) SuffixIf(sql safeString, include bool, args ...interface{}) deleteBuilder {
	return b.SuffixExprIf(Expr(sql, args...), include)
}

// SuffixExpr adds an expression to the end of the query
func (b deleteBuilder) SuffixExprIf(expr Sqlizer, include bool) deleteBuilder {
	if include {
		return b.SuffixExpr(expr)
	}
	return b
}
