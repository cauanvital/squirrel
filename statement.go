package squirrel2

// StatementBuilderType is the type of StatementBuilder.
type statementBuilderType struct {
	placeholderFormat PlaceholderFormat
	runWith           BaseRunner
	whereParts        []Sqlizer
}

func StatementBuilderType() statementBuilderType {
	return statementBuilderType{
		placeholderFormat: Question,
		whereParts:        make([]Sqlizer, 0),
	}
}

// Select returns a SelectBuilder for this StatementBuilderType.
func (b statementBuilderType) Select(columns ...safeString) selectBuilder {
	return SelectBuilder(b).Columns(columns...)
}

// Insert returns a InsertBuilder for this StatementBuilderType.
func (b statementBuilderType) Insert(into safeString) insertBuilder {
	return InsertBuilder(b).Into(into)
}

// Replace returns a InsertBuilder for this StatementBuilderType with the
// statement keyword set to "REPLACE".
func (b statementBuilderType) Replace(into safeString) insertBuilder {
	return InsertBuilder(b).statementKeyword("REPLACE").Into(into)
}

// Update returns a UpdateBuilder for this StatementBuilderType.
func (b statementBuilderType) Update(table safeString) updateBuilder {
	return UpdateBuilder(b).Table(table)
}

// Delete returns a DeleteBuilder for this StatementBuilderType.
func (b statementBuilderType) Delete(from safeString) deleteBuilder {
	return DeleteBuilder(b).From(from)
}

// PlaceholderFormat sets the PlaceholderFormat field for any child builders.
func (b statementBuilderType) PlaceholderFormat(f PlaceholderFormat) statementBuilderType {
	b.placeholderFormat = f
	return b
}

// RunWith sets the RunWith field for any child builders.
func (b statementBuilderType) RunWith(runner BaseRunner) statementBuilderType {
	switch r := runner.(type) {
	case StdSqlCtx:
		runner = WrapStdSqlCtx(r)
	case StdSql:
		runner = WrapStdSql(r)
	}
	b.runWith = runner
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b statementBuilderType) Where(expr Sqlizer) statementBuilderType {
	b.whereParts = append(b.whereParts, expr)
	return b
}

// StatementBuilder is a parent builder for other builders, e.g. SelectBuilder.
var StatementBuilder = StatementBuilderType()

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...safeString) selectBuilder {
	return StatementBuilder.Select(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into safeString) insertBuilder {
	return StatementBuilder.Insert(into)
}

// Replace returns a new InsertBuilder with the statement keyword set to
// "REPLACE" and with the given table name.
//
// See InsertBuilder.Into.
func Replace(into safeString) insertBuilder {
	return StatementBuilder.Replace(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table safeString) updateBuilder {
	return StatementBuilder.Update(table)
}

// Delete returns a new DeleteBuilder with the given table name.
//
// See DeleteBuilder.Table.
func Delete(from safeString) deleteBuilder {
	return StatementBuilder.Delete(from)
}

// Case returns a new CaseBuilder
// "what" represents optional case value
func Case(what ...Sqlizer) caseBuilder {
	b := CaseBuilder()
	switch len(what) {
	case 0:
		return b
	default:
		return b.what(what[0])
	}
}
