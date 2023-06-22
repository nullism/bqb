package bqb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Dialect holds the Query dialect
type Dialect string

const (
	// PGSQL postgres dialect
	PGSQL Dialect = "postgres"
	// MYSQL MySQL dialect
	MYSQL Dialect = "mysql"
	// RAW dialect uses no parameter conversion
	RAW Dialect = "raw"
	// SQL generic dialect
	SQL Dialect = "sql"

	paramPh = "{{xX_PARAM_Xx}}"
)

// JsonMap is a custom type which tells bqb to convert the parameter to
// a JSON object without requiring reflection.
type JsonMap map[string]interface{}

// JsonList is a type that tells bqb to convert the parameter to a JSON
// list without requiring reflection.
type JsonList []interface{}

// Identifiers is a type that tells bqb to quote the prameter and inline it rather than using query
// parameters. This allows for (somewhat)safely parameterizing table & column names.
type Identifiers []string

func dialectReplace(dialect Dialect, sql string, params []interface{}) (string, []interface{}, error) {
	newParams := make([]interface{}, 0, len(params))
	for i, param := range params {
		switch v := param.(type) {
		case Identifiers:
			sql = strings.Replace(sql, paramPh, quoteIdentifiers(v, dialect), 1)
		default:
			newParams = append(newParams, param)
			if dialect == RAW {
				p, err := paramToRaw(param)
				if err != nil {
					return "", nil, err
				}
				sql = strings.Replace(sql, paramPh, p, 1)
			} else if dialect == PGSQL {
				sql = strings.ReplaceAll(sql, "??", "?")
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
			} else {
				sql = strings.Replace(sql, paramPh, "?", 1)
			}
		}
	}
	return sql, newParams, nil
}

func makePart(text string, args ...interface{}) QueryPart {
	tempPh := "XXX___XXX"
	originalText := text
	text = strings.ReplaceAll(text, "??", tempPh)

	var newArgs []interface{}
	errs := make([]error, 0)

	for _, arg := range args {
		switch v := arg.(type) {

		case driver.Valuer:
			text = strings.Replace(text, "?", paramPh, 1)
			val, err := v.Value()
			if err != nil {
				errs = append(errs, err)
			} else {
				newArgs = append(newArgs, val)
			}
		case []int:
			newPh := []string{}
			for _, i := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, i)
			}
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

		case []*int:
			newPh := []string{}
			for _, i := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, i)
			}
			if len(newPh) > 0 {
				text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)
			} else {
				text = strings.Replace(text, "?", paramPh, 1)
				newArgs = append(newArgs, nil)
			}

		case []string:
			newPh := []string{}
			for _, s := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, s)
			}
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

		case []*string:
			newPh := []string{}
			for _, s := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, s)
			}
			if len(newPh) > 0 {
				text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)
			} else {
				text = strings.Replace(text, "?", paramPh, 1)
				newArgs = append(newArgs, nil)
			}

		case []interface{}:
			newPh := []string{}
			for _, s := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, s)
			}
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

		case *Query:
			if v == nil {
				text = strings.Replace(text, "?", paramPh, 1)
				newArgs = append(newArgs, nil)
				continue
			}
			sql, params, _ := v.toSql()
			text = strings.Replace(text, "?", sql, 1)
			newArgs = append(newArgs, params...)

		case JsonMap, JsonList:
			bytes, err := json.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("cann jsonify struct: %v", err))
			}
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, string(bytes))

		case *JsonMap, *JsonList:
			bytes, err := json.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("cann jsonify struct: %v", err))
			}
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, string(bytes))

		case Identifiers:
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, v)

		default:
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, v)
		}
	}
	extraCount := strings.Count(text, "?")
	if extraCount > 0 {
		panic(fmt.Sprintf("extra ? in text: %v (%d args)", originalText, len(newArgs)))
	}

	paramCount := strings.Count(text, paramPh)
	if paramCount < len(newArgs) {
		panic(fmt.Sprintf("missing ? in text: %v (%d args)", originalText, len(newArgs)))
	}

	text = strings.ReplaceAll(text, tempPh, "??")

	return QueryPart{
		Text:   text,
		Params: newArgs,
		Errs:   errs,
	}
}

func paramToRaw(param interface{}) (string, error) {
	switch p := param.(type) {
	case bool:
		return fmt.Sprintf("%v", p), nil
	case float32, float64, int, int8, int16, int32, int64,
		uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", p), nil
	case *int:
		if p == nil {
			return "NULL", nil
		}
		return fmt.Sprintf("%v", *p), nil
	case string:
		return fmt.Sprintf("'%v'", p), nil
	case *string:
		if p == nil {
			return "NULL", nil
		}
		return fmt.Sprintf("'%v'", *p), nil
	case nil:
		return "NULL", nil
	default:
		return "", fmt.Errorf("unsupported type for Raw query: %T", p)
	}
}

func quoteIdentifiers(names Identifiers, dialect Dialect) string {
	qChar := `"`
	if dialect == MYSQL {
		qChar = "`"
	}
	quoted := make([]string, len(names))
	for i, name := range names {
		quoted[i] = qChar + strings.Replace(name, qChar, qChar+qChar, -1) + qChar
	}
	return strings.Join(quoted, ".")
}
