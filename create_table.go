package bqb

import "strings"

type createTable struct {
	dialect string
	table   Expr
	cols    []Expr
	selectQ *selectQ
}

func CreateTable(intf interface{}) *createTable {
	return &createTable{
		dialect: SQL,
		table:   intfToExpr(intf),
	}
}

func (g *createTable) Postgres() *createTable {
	g.dialect = PGSQL
	return g
}

func (g *createTable) Mysql() *createTable {
	g.dialect = MYSQL
	return g
}

func (g *createTable) Raw() *createTable {
	g.dialect = RAW
	return g
}

func (c *createTable) Cols(intfs ...interface{}) *createTable {
	c.cols = append(c.cols, getExprs(intfs)...)
	return c
}

func (c *createTable) Select(s *selectQ) *createTable {
	c.selectQ = s
	return c
}

func (c *createTable) ToSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	sql += "CREATE TABLE "
	sql += c.table.F + " "
	params = append(params, c.table.V...)

	if len(c.cols) > 0 {
		sql += "("
		nsql, nparams := exprsToSql(c.cols)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += ") "
	}

	if c.selectQ != nil {
		sql += "AS "
		qs, qp, err := c.selectQ.toSql()
		if err != nil {
			return "", nil, err
		}
		sql += qs
		params = append(params, qp...)
	}

	sql = dialectReplace(c.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
