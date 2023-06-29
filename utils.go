package bqb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

func dialectReplace(dialect Dialect, sql string, params []any) (string, error) {
	if dialect == MYSQL || dialect == SQL {
		sql = strings.ReplaceAll(sql, paramPh, "?")
	}
	for i, param := range params {
		if dialect == RAW {
			p, err := paramToRaw(param)
			if err != nil {
				return "", err
			}
			sql = strings.Replace(sql, paramPh, p, 1)
		} else if dialect == PGSQL {
			sql = strings.ReplaceAll(sql, "??", "?")
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}
	}
	return sql, nil
}

func convertArg(text string, arg any) (string, []any, []error) {
	var newArgs []any
	var errs []error

	switch v := arg.(type) {

	case Embedder:
		text = strings.Replace(text, "?", v.RawValue(), 1)

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

	case []any:
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
			return text, newArgs, errs
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

	case Embedded:
		text = strings.Replace(text, "?", string(v), 1)

	default:
		text = strings.Replace(text, "?", paramPh, 1)
		newArgs = append(newArgs, v)
	}

	return text, newArgs, errs
}

func makePart(text string, args ...any) QueryPart {
	tempPh := "XXX___XXX"
	originalText := text
	text = strings.ReplaceAll(text, "??", tempPh)

	var newArgs []any
	errs := make([]error, 0)

	for _, arg := range args {
		argText, fArgs, argErrs := convertArg(text, arg)
		if len(argErrs) > 0 {
			errs = append(errs, argErrs...)
		}
		newArgs = append(newArgs, fArgs...)
		text = argText
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

func paramToRaw(param any) (string, error) {
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
