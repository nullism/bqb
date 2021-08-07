package bqb

import (
	"fmt"
	"strings"
)

const (
	PGSQL = "postgres"
	MYSQL = "mysql"
	RAW   = "raw"
	SQL   = "sql"

	paramPh = "{{xX_PARAM_Xx}}"
)

func paramToRaw(param interface{}) string {
	switch p := param.(type) {
	case int:
		return fmt.Sprintf("%v", p)
	case string:
		return fmt.Sprintf("'%v'", p)
	case nil:
		return "NULL"
	default:
		panic(fmt.Sprintf("cannot convert type %T", p))
	}
}

func dialectReplace(dialect string, sql string, params []interface{}) string {
	for i, param := range params {
		if dialect == RAW {
			sql = strings.Replace(sql, paramPh, paramToRaw(param), 1)
		} else if dialect == MYSQL || dialect == SQL {
			sql = strings.Replace(sql, paramPh, "?", 1)
		} else if dialect == PGSQL {
			sql = strings.ReplaceAll(sql, "??", "?")
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}
	}
	return sql
}
