package squirrel

// SelectIf returns a SelectBuilder for this StatementBuilderType, only with columns with true boolean value.
func (b StatementBuilderType) SelectIf(columns ...valIf[safeString]) selectBuilder {
	cols := make([]safeString, 0)
	for _, v := range columns {
		if v.Include {
			cols = append(cols, v.Value)
		}
	}
	return SelectBuilder().Columns(cols...)
}

// SelectIf returns a new SelectBuilder, optionally setting some result columns, only with columns with true boolean value.
//
// See SelectBuilder.Columns.
func SelectIf(columns ...valIf[safeString]) selectBuilder {
	cols := make([]safeString, 0)
	for _, v := range columns {
		if v.Include {
			cols = append(cols, v.Value)
		}
	}
	return StatementBuilder.Select(cols...)
}
