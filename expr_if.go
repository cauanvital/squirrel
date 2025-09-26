package squirrel

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

type exprIf struct {
	expr
	include bool
}

// ExprIf creates an expression from a SQL fragment and arguments that will be written only if 'include' is true.
//
// Ex:
//
//	ExprIf(Expr("FROM_UNIXTIME(?)", t), true)
func ExprIf(e expr, include bool) Sqlizer {
	return exprIf{expr: e, include: include}
}

func (eIf exprIf) ToSql() (sql string, args []any, err error) {
	if !eIf.include {
		return "", nil, nil
	}
	e := eIf.expr

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

// EqIf is syntactic sugar for use with Where/Having/Set methods that will be written only if 'include' is true.
//
// Ex:
//
//	.Where(EqIf{{"id": 1}, true}) == "id = 1"
type EqIf struct {
	Eq
	include bool
}

func (eqIf EqIf) toSQL(useNotOpr bool) (sql string, args []interface{}, err error) {
	if !eqIf.include {
		return "", nil, nil
	}
	eq := eqIf.Eq

	if len(eq) == 0 {
		// Empty Sql{} evaluates to true.
		sql = sqlTrue
		return
	}

	var (
		exprs       []string
		equalOpr    = "="
		inOpr       = "IN"
		nullOpr     = "IS"
		inEmptyExpr = sqlFalse
	)

	if useNotOpr {
		equalOpr = "<>"
		inOpr = "NOT IN"
		nullOpr = "IS NOT"
		inEmptyExpr = sqlTrue
	}

	sortedKeys := getSortedKeys(eq)
	for _, key := range sortedKeys {
		var expr string
		val := eq[key]

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		r := reflect.ValueOf(val)
		if r.Kind() == reflect.Ptr {
			if r.IsNil() {
				val = nil
			} else {
				val = r.Elem().Interface()
			}
		}

		if val == nil {
			expr = fmt.Sprintf("%s %s NULL", key, nullOpr)
		} else {
			if isListType(val) {
				valVal := reflect.ValueOf(val)
				if valVal.Len() == 0 {
					expr = inEmptyExpr
					if args == nil {
						args = []interface{}{}
					}
				} else {
					for i := 0; i < valVal.Len(); i++ {
						args = append(args, valVal.Index(i).Interface())
					}
					expr = fmt.Sprintf("%s %s (%s)", key, inOpr, Placeholders(valVal.Len()))
				}
			} else {
				expr = fmt.Sprintf("%s %s ?", key, equalOpr)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

func (eq EqIf) ToSql() (sql string, args []interface{}, err error) {
	return eq.toSQL(false)
}

// NotEqIf is syntactic sugar for use with Where/Having/Set methods that will be written only if 'include' is true.
// Ex:
//
//	.Where(NotEq{"id": 1}) == "id <> 1"
type NotEqIf EqIf

func (neqIf NotEqIf) ToSql() (sql string, args []interface{}, err error) {
	return EqIf(neqIf).toSQL(true)
}

// Like is syntactic sugar for use with LIKE conditions.
// Ex:
//
//	.Where(Like{"name": "%irrel"})
type LikeIf struct {
	Like
	include bool
}

func (lkIf LikeIf) toSql(opr string) (sql string, args []interface{}, err error) {
	if !lkIf.include {
		return "", nil, nil
	}
	lk := lkIf.Like

	var exprs []string
	for key, val := range lk {
		expr := ""

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		if val == nil {
			err = fmt.Errorf("cannot use null with like operators")
			return
		} else {
			if isListType(val) {
				err = fmt.Errorf("cannot use array or slice with like operators")
				return
			} else {
				expr = fmt.Sprintf("%s %s ?", key, opr)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

func (lk LikeIf) ToSql() (sql string, args []interface{}, err error) {
	return lk.toSql("LIKE")
}

// NotLike is syntactic sugar for use with LIKE conditions.
// Ex:
//
//	.Where(NotLike{"name": "%irrel"})
type NotLikeIf LikeIf

func (nlkIf NotLikeIf) ToSql() (sql string, args []interface{}, err error) {
	return LikeIf(nlkIf).toSql("NOT LIKE")
}

// ILike is syntactic sugar for use with ILIKE conditions.
// Ex:
//
//	.Where(ILike{"name": "sq%"})
type ILikeIf LikeIf

func (ilkIf ILikeIf) ToSql() (sql string, args []interface{}, err error) {
	return LikeIf(ilkIf).toSql("ILIKE")
}

// NotILike is syntactic sugar for use with ILIKE conditions.
// Ex:
//
//	.Where(NotILike{"name": "sq%"})
type NotILikeIf LikeIf

func (nilkIf NotILikeIf) ToSql() (sql string, args []interface{}, err error) {
	return LikeIf(nilkIf).toSql("NOT ILIKE")
}

// Lt is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//
//	.Where(Lt{"id": 1})
type LtIf struct {
	Lt
	include bool
}

func (ltIf LtIf) toSql(opposite, orEq bool) (sql string, args []interface{}, err error) {
	if !ltIf.include {
		return "", nil, nil
	}
	lt := ltIf.Lt

	var (
		exprs []string
		opr   = "<"
	)

	if opposite {
		opr = ">"
	}

	if orEq {
		opr = fmt.Sprintf("%s%s", opr, "=")
	}

	sortedKeys := getSortedKeys(lt)
	for _, key := range sortedKeys {
		var expr string
		val := lt[key]

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		if val == nil {
			err = fmt.Errorf("cannot use null with less than or greater than operators")
			return
		}
		if isListType(val) {
			err = fmt.Errorf("cannot use array or slice with less than or greater than operators")
			return
		}
		expr = fmt.Sprintf("%s %s ?", key, opr)
		args = append(args, val)

		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

func (ltIf LtIf) ToSql() (sql string, args []interface{}, err error) {
	return ltIf.toSql(false, false)
}

// LtOrEq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//
//	.Where(LtOrEq{"id": 1}) == "id <= 1"
type LtOrEqIf LtIf

func (ltOrEqIf LtOrEqIf) ToSql() (sql string, args []interface{}, err error) {
	return LtIf(ltOrEqIf).toSql(false, true)
}

// Gt is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//
//	.Where(Gt{"id": 1}) == "id > 1"
type GtIf LtIf

func (gtIf GtIf) ToSql() (sql string, args []interface{}, err error) {
	return LtIf(gtIf).toSql(true, false)
}

// GtOrEq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//
//	.Where(GtOrEq{"id": 1}) == "id >= 1"
type GtOrEqIf LtIf

func (gtOrEqIf GtOrEqIf) ToSql() (sql string, args []interface{}, err error) {
	return LtIf(gtOrEqIf).toSql(true, true)
}
