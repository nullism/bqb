package bqb

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Dialect holds the Query dialect
type Dialect string

const (
	// PostGres dialect
	PGSQL Dialect = "postgres"
	// MySQL dialect
	MYSQL Dialect = "mysql"
	// Raw dialect uses no parameter conversion
	RAW Dialect = "raw"
	// Generic SQL dialect
	SQL Dialect = "sql"

	paramPh = "{{xX_PARAM_Xx}}"
)

// JsonMap is a custom type which tells bqb to convert the parameter to
// a JSON object without requiring reflection.
type JsonMap map[string]interface{}

// JsonList is a type that tells bqb to convert the parameter to a JSON
// list without requiring reflection.
type JsonList []interface{}

func dialectReplace(dialect Dialect, sql string, params []interface{}) (string, error) {
	for i, param := range params {
		if dialect == RAW {
			p, err := paramToRaw(param)
			if err != nil {
				return "", err
			}
			sql = strings.Replace(sql, paramPh, p, 1)
		} else if dialect == MYSQL || dialect == SQL {
			sql = strings.Replace(sql, paramPh, "?", 1)
		} else if dialect == PGSQL {
			sql = strings.ReplaceAll(sql, "??", "?")
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}
	}
	return sql, nil
}

func makePart(text string, args ...interface{}) QueryPart {
	tempPh := "XXX___XXX"
	originalText := text
	text = strings.ReplaceAll(text, "??", tempPh)

	var newArgs []interface{}

	for _, arg := range args {
		switch v := arg.(type) {

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
