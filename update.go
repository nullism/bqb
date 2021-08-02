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
}

func Update(exprs ...interface{}) *update {
	return &update{
		dialect: SQL,
		update:  getExprs(exprs),
	}
}

func UpdateSql() *update {
	return &update{dialect: SQL}
}

func UpdateMysql() *update {
	return &update{dialect: MYSQL}
}

func UpdateRaw() *update {
	return &update{dialect: PGSQL}
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

	sql = dialectReplace(u.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
