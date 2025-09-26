package squirrel

import (
	"bytes"
	"strings"
)

type exprIf struct {
	expression expr
	include    bool
}

// ExprIf is a Sqlizer that conditionally wraps an expression.
//
// Ex:
//
//	ExprIf(Expr("FROM_UNIXTIME(?)", t), true) == "FROM_UNIXTIME(t)"
//	ExprIf(Expr("FROM_UNIXTIME(?)", t), false) == ""
func ExprIf(expression expr, include bool) Sqlizer {
	return exprIf{expression: expression, include: include}
}

func (eIf exprIf) ToSql() (sql string, args []any, err error) {
	if !eIf.include {
		return "", nil, nil
	}
	e := eIf.expression

	simple := true
	for _, arg := range e.args {
		if _, ok := arg.(Sqlizer); ok {
			simple = false
		}
	}
	if simple {
		return string(e.sql), e.args, nil
	}

	buf := &bytes.Buffer{}
	ap := e.args
	sp := string(e.sql)

	var isql string
	var iargs []any

	for err == nil && len(ap) > 0 && len(sp) > 0 {
		i := strings.Index(sp, "?")
		if i < 0 {
			// no more placeholders
			break
		}
		if len(sp) > i+1 && sp[i+1:i+2] == "?" {
			// escaped "??"; append it and step past
			buf.WriteString(sp[:i+2])
			sp = sp[i+2:]
			continue
		}

		if as, ok := ap[0].(Sqlizer); ok {
			// sqlizer argument; expand it and append the result
			isql, iargs, err = as.ToSql()
			buf.WriteString(sp[:i])
			buf.WriteString(isql)
			args = append(args, iargs...)
		} else {
			// normal argument; append it and the placeholder
			buf.WriteString(sp[:i+1])
			args = append(args, ap[0])
		}

		// step past the argument and placeholder
		ap = ap[1:]
		sp = sp[i+1:]
	}

	// append the remaining sql and arguments
	buf.WriteString(sp)
	return buf.String(), append(args, ap...), err
}

// SqIf is a Sqlizer that conditionally wraps another Sqlizer.
//
// Ex:
//
//	.Where(SqIf{Eq{"id": 1}, true}) == "id = 1"
//	.Where(SqIf{Eq{"id": 2}, false}) == ""
type SqIf struct {
	clause  Sqlizer
	include bool
}

func (sqIf SqIf) ToSql() (sql string, args []interface{}, err error) {
	if sqIf.include {
		return sqIf.clause.ToSql()
	}
	return "", nil, nil
}
