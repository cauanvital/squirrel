package squirrel2

// WriteSql converts Sqlizer to SQL strings and writes it to buffer
func (b *sqlizerBuffer) WriteSqlIf(item Sqlizer, include bool) {
	if include {
		b.WriteSql(item)
	}
}

// When adds "WHEN ... THEN ..." part to CASE construct
func (b caseBuilder) WhenIf(when Sqlizer, then Sqlizer, include bool) caseBuilder {
	if include {
		return b.When(when, then)
	}
	return b
}

// What sets optional "ELSE ..." part for CASE construct
func (b caseBuilder) ElseIf(expr Sqlizer, include bool) caseBuilder {
	if include {
		return b.Else(expr)
	}
	return b
}
