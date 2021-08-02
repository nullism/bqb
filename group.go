package bqb

import "strings"

type group struct {
	dialect string
	groups  []Expr
}

func Group(exprs ...interface{}) *group {
	return &group{
		dialect: SQL,
		groups:  getExprs(exprs),
	}
}

func (g *group) Postgres() *group {
	g.dialect = PGSQL
	return g
}

func (g *group) Mysql() *group {
	g.dialect = MYSQL
	return g
}

func (g *group) Raw() *group {
	g.dialect = RAW
	return g
}

func (g *group) ToSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	if len(g.groups) > 0 {
		sql += ""
		nsql, nparams := exprsToSql(g.groups)
		sql += strings.Join(nsql, " ")
		params = append(params, nparams...)
		sql += " "
	}

	sql = dialectReplace(g.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
