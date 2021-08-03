package bqb

import (
	"fmt"
	"strings"
)

type update struct {
	dialect string
	update  []Expr
	set     []Expr
	where   []Expr
	offset  int
	limit   int
}

func Update(exprs ...interface{}) *update {
	return &update{
		dialect: SQL,
		update:  getExprs(exprs),
	}
}

func (u *update) Postgres() *update {
	u.dialect = PGSQL
	return u
}

func (u *update) Mysql() *update {
	u.dialect = MYSQL
	return u
}

func (u *update) Raw() *update {
	u.dialect = RAW
	return u
}

func (u *update) Set(exprs ...interface{}) *update {
	newExprs := getExprs(exprs)
	u.set = append(u.set, newExprs...)
	return u
}

func (u *update) Where(exprs ...interface{}) *update {
	newExprs := getExprs(exprs)
	u.where = append(u.where, newExprs...)
	return u
}

func (u *update) Offset(offset int) *update {
	u.offset = offset
	return u
}

func (u *update) Limit(limit int) *update {
	u.limit = limit
	return u
}

func (u *update) Print() {
	sql, params, err := u.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (u *update) ToSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	if len(u.update) > 0 {
		sql += "UPDATE "
		nsql, nparams := exprsToSql(u.update)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(u.set) > 0 {
		sql += "SET "
		nsql, nparams := exprsToSql(u.set)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(u.where) > 0 {
		sql += "WHERE "
		nsql, nparams := exprsToSql(u.where)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if u.offset != 0 {
		sql += fmt.Sprintf("OFFSET %d ", u.offset)
	}

	if u.limit != 0 {
		sql += fmt.Sprintf("LIMIT %d ", u.limit)
	}

	sql = dialectReplace(u.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
