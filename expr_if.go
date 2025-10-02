package squirrel2

type exprIf struct {
	expression Sqlizer
	include    bool
}

// ExprIf is a Sqlizer that conditionally wraps an expression.
//
// Ex:
//
//	ExprIf(Expr("FROM_UNIXTIME(?)", t), true) == "FROM_UNIXTIME(t)"
//	ExprIf(Expr("FROM_UNIXTIME(?)", t), false) == ""
func ExprIf(expression Sqlizer, include bool) Sqlizer {
	return exprIf{expression: expression, include: include}
}

func (eIf exprIf) ToSql() (sql string, args []interface{}, err error) {
	if eIf.include {
		return eIf.expression.ToSql()
	}
	return "", nil, nil
}
