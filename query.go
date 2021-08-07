package bqb

import (
	"fmt"
	"strings"
)

type part struct {
	Text   string
	Params []interface{}
}

type Query struct {
	dialect string
	Parts   []part
	Prepend string
}

func makePart(text string, args ...interface{}) part {
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

		case []string:
			newPh := []string{}
			for _, s := range v {
				newPh = append(newPh, paramPh)
				newArgs = append(newArgs, s)
			}
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

		case *Query:
			sql, params, _ := v.toSql()
			text = strings.Replace(text, "?", sql, 1)
			newArgs = append(newArgs, params...)

		default:
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, v)
		}
	}
	extraCount := strings.Count(text, "?")
	if extraCount > 0 {
		panic(fmt.Sprintf("extra ? in text: %v", originalText))
	}

	paramCount := strings.Count(text, paramPh)
	if paramCount < len(newArgs) {
		panic(fmt.Sprintf("missing ? in text: %v", originalText))
	}

	text = strings.ReplaceAll(text, tempPh, "??")

	return part{
		Text:   text,
		Params: newArgs,
	}
}

func New(text string, args ...interface{}) *Query {
	q := &Query{
		dialect: SQL,
	}
	q.Parts = append(q.Parts, makePart(text, args...))
	return q
}

func Empty(prep ...string) *Query {
	return &Query{
		Prepend: strings.Join(prep, " "),
	}
}

func (q *Query) Space(text string, args ...interface{}) *Query {
	return q.Join(" ", text, args...)
}

func (q *Query) And(text string, args ...interface{}) *Query {
	return q.Join(" AND ", text, args...)
}

func (q *Query) Or(text string, args ...interface{}) *Query {
	return q.Join(" OR ", text, args...)
}

func (q *Query) Comma(text string, args ...interface{}) *Query {
	return q.Join(",", text, args...)
}

func (q *Query) Concat(text string, args ...interface{}) *Query {
	return q.Join("", text, args...)
}

func (q *Query) Join(sep, text string, args ...interface{}) *Query {
	if len(q.Parts) > 0 {
		q.Parts = append(q.Parts, makePart(sep+text, args...))
	} else {
		q.Parts = append(q.Parts, makePart(text, args...))
	}

	return q
}

func (q *Query) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	fmt.Printf("PARAMS: %v\n", params)
	fmt.Printf("ERROR: %v\n", err)
}

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}

	return dialectReplace(q.dialect, sql, params), params, nil
}

func (q *Query) ToPsql() (string, []interface{}, error) {
	q.dialect = PGSQL
	return q.ToSql()
}

func (q *Query) ToRaw() (string, error) {
	q.dialect = RAW
	sql, _, err := q.ToSql()
	return sql, err
}

func (q *Query) toSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	if q.Prepend != "" && len(q.Parts) > 0 {
		sql = q.Prepend + " "
	}

	for _, p := range q.Parts {
		sql += p.Text
		params = append(params, p.Params...)
	}

	return strings.TrimSpace(sql), params, nil
}
