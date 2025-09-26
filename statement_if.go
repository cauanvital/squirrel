package squirrel

import "github.com/lann/builder"

type ValIf[T any] struct {
	Value   T
	Include bool
}

// SelectIf returns a SelectBuilder for this StatementBuilderType, only with columns with true boolean value.
func (b StatementBuilderType) SelectIf(columnsIf ...ValIf[string]) SelectBuilder {
	columns := make([]string, 0)
	for _, v := range columnsIf {
		if v.Include {
			columns = append(columns, v.Value)
		}
	}
	return SelectBuilder(b).Columns(columns...)
}

// WhereIf adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b StatementBuilderType) WhereIf(pred interface{}, args ...interface{}) StatementBuilderType {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(StatementBuilderType)
}

// SelectIf returns a new SelectBuilder, optionally setting some result columns, only with columns with true boolean value.
//
// See SelectBuilder.Columns.
func SelectIf(columnsIf ...ValIf[string]) SelectBuilder {
	columns := make([]string, 0)
	for _, v := range columnsIf {
		if v.Include {
			columns = append(columns, v.Value)
		}
	}
	return StatementBuilder.Select(columns...)
}
