package bqb

type Insert struct {
	dialect string
	into    []Expr
	cols    []Expr
	vals    []Expr
}

func InsertPsql() *Insert {
	return &Insert{dialect: PGSQL}
}

func (i *Insert) Into(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.into = append(i.into, newExprs...)
	return i
}

func (i *Insert) ToSql() {
	sql := "INSERT "

	if len(i.into) > 0 {
		sql += "X"
	}
}
