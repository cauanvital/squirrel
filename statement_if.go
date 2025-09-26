package squirrel

// SelectIf returns a SelectBuilder for this StatementBuilderType, only with columns with true boolean value.
func (b StatementBuilderType) SelectIf(columnsIf ...ValIf[safeString]) SelectBuilder {
	columns := make([]safeString, 0)
	for _, v := range columnsIf {
		if v.Include {
			columns = append(columns, v.Value)
		}
	}
	return SelectBuilder(b).Columns(columns...)
}

// SelectIf returns a new SelectBuilder, optionally setting some result columns, only with columns with true boolean value.
//
// See SelectBuilder.Columns.
func SelectIf(columnsIf ...ValIf[safeString]) SelectBuilder {
	columns := make([]safeString, 0)
	for _, v := range columnsIf {
		if v.Include {
			columns = append(columns, v.Value)
		}
	}
	return StatementBuilder.Select(columns...)
}
