package squirrel

import (
	"bytes"
	"errors"
)

// sqlizerBuffer is a helper that allows to write many Sqlizers one by one
// without constant checks for errors that may come from Sqlizer
type sqlizerBuffer struct {
	bytes.Buffer
	args []interface{}
	err  error
}

// WriteSql converts Sqlizer to SQL strings and writes it to buffer
func (b *sqlizerBuffer) WriteSql(item Sqlizer) {
	if b.err != nil {
		return
	}

	var str string
	var args []interface{}
	str, args, b.err = nestedToSql(item)

	if b.err != nil {
		return
	}

	b.WriteString(str)
	b.WriteByte(' ')
	b.args = append(b.args, args...)
}

func (b *sqlizerBuffer) ToSql() (string, []interface{}, error) {
	return b.String(), b.args, b.err
}

// whenPart is a helper structure to describe SQLs "WHEN ... THEN ..." expression
type whenPart struct {
	when Sqlizer
	then Sqlizer
}

// caseData holds all the data required to build a CASE SQL construct
type caseData struct {
	What      Sqlizer
	WhenParts []whenPart
	Else      Sqlizer
}

// ToSql implements Sqlizer
func (d *caseData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.WhenParts) == 0 {
		err = errors.New("case expression must contain at lease one WHEN clause")
		return
	}

	sql := sqlizerBuffer{}

	sql.WriteString("CASE ")
	if d.What != nil {
		sql.WriteSql(d.What)
	}

	for _, p := range d.WhenParts {
		sql.WriteString("WHEN ")
		sql.WriteSql(p.when)
		sql.WriteString("THEN ")
		sql.WriteSql(p.then)
	}

	if d.Else != nil {
		sql.WriteString("ELSE ")
		sql.WriteSql(d.Else)
	}
	sql.WriteString("END")

	return sql.ToSql()
}

// CaseBuilder builds SQL CASE construct which could be used as parts of queries.
type caseBuilder struct {
	data caseData
}

func CaseBuilder() caseBuilder {
	return caseBuilder{
		data: caseData{
			WhenParts: make([]whenPart, 0),
		},
	}
}

// ToSql builds the query into a SQL string and bound args.
func (b caseBuilder) ToSql() (string, []interface{}, error) {
	return b.data.ToSql()
}

// MustSql builds the query into a SQL string and bound args.
// It panics if there are any errors.
func (b caseBuilder) MustSql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// what sets the optional value for the CASE construct, e.g., "CASE status ..."
func (b caseBuilder) what(expr Sqlizer) caseBuilder {
	b.data.What = expr
	return b
}

// When adds "WHEN ... THEN ..." part to CASE construct
func (b caseBuilder) When(when Sqlizer, then Sqlizer) caseBuilder {
	// TODO: performance hint: replace slice of WhenPart with just slice of parts
	// where even indices of the slice belong to "when"s and odd indices belong to "then"s
	b.data.WhenParts = append(b.data.WhenParts, whenPart{when, then})
	return b
}

// Else sets the optional "ELSE ..." part for the CASE construct.
func (b caseBuilder) Else(expr Sqlizer) caseBuilder {
	b.data.Else = expr
	return b
}
