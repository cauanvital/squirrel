package squirrel2

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b insertBuilder) PlaceholderFormatIf(f PlaceholderFormat, include bool) insertBuilder {
	if include {
		return b.PlaceholderFormat(f)
	}
	return b
}

// Runner methods

// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
func (b insertBuilder) RunWithIf(runner BaseRunner, include bool) insertBuilder {
	if include {
		return b.RunWith(runner)
	}
	return b
}

// Prefix adds an expression to the beginning of the query
func (b insertBuilder) PrefixIf(sql safeString, include bool, args ...interface{}) insertBuilder {
	if include {
		return b.PrefixExpr(Expr(sql, args...))
	}
	return b
}

// PrefixExpr adds an expression to the very beginning of the query
func (b insertBuilder) PrefixExprIf(expr Sqlizer, include bool) insertBuilder {
	if include {
		return b.PrefixExpr(expr)
	}
	return b
}

// Options adds keyword options before the INTO clause of the query.
func (b insertBuilder) OptionsIf(options ...valIf[safeString]) insertBuilder {
	for _, v := range options {
		b = b.OptionIf(v.Value, v.Include)
	}
	return b
}

// OptionIf adds a single select option to the query if include is true.
func (b insertBuilder) OptionIf(option safeString, include bool) insertBuilder {
	if include {
		return b.Options(option)
	}
	return b
}

// Columns adds insert columns to the query.
func (b insertBuilder) ColumnsIf(columns ...valIf[safeString]) insertBuilder {
	for _, v := range columns {
		b = b.ColumnIf(v.Value, v.Include)
	}
	return b
}

func (b insertBuilder) ColumnIf(column safeString, include bool) insertBuilder {
	if include {
		return b.Columns(column)
	}
	return b
}

// Values adds a single row's values to the query.
func (b insertBuilder) ValuesIf(values ...valIf[interface{}]) insertBuilder {
	for _, v := range values {
		b = b.ValueIf(v.Value, v.Include)
	}
	return b
}

func (b insertBuilder) ValueIf(value interface{}, include bool) insertBuilder {
	if include {
		return b.Values(value)
	}
	return b
}

// Suffix adds an expression to the end of the query
func (b insertBuilder) SuffixIf(sql safeString, include bool, args ...interface{}) insertBuilder {
	if include {
		return b.Suffix(sql, args...)
	}
	return b
}

// SuffixExpr adds an expression to the end of the query
func (b insertBuilder) SuffixExprIf(expr Sqlizer, include bool) insertBuilder {
	if include {
		return b.SuffixExpr(expr)
	}
	return b
}

// SetMap set columns and values for insert builder from a map of column name and value
// note that it will reset all previous columns and values was set if any
func (b insertBuilder) SetMapIf(clauses map[safeString]interface{}, include bool) insertBuilder {
	if include {
		return b.SetMap(clauses)
	}
	return b
}

// Select set Select clause for insert query
// If Values and Select are used, then Select has higher priority
func (b insertBuilder) SelectIf(sb selectBuilder, include bool) insertBuilder {
	if include {
		return b.Select(sb)
	}
	return b
}
