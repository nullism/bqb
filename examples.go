package main

import (
	"fmt"

	"github.com/nullism/bqb"
)

func main() {

	q := bqb.Query{
		S: "*, t.name, t.id",
		F: "my_table t",
		J: []string{"my_other_tab ot ON ot.id = t.id"},
		W: [][]bqb.Expr{
			{
				{
					F: "ST_Distance(geom, ?) < ?",
					V: []interface{}{"foo", 23},
				},
				{
					F: "t.count < ?",
					V: []interface{}{1000},
				},
			},
			{
				{
					F: "ST_Dwithin(geom, GeomFromEWKT(?))",
					V: []interface{}{"SRID=4326;POINT(44 -111)"},
				},
			},
			{
				{
					F: "t.id IN (SELECT id FROM other_table WHERE id like t.id AND count > ?)",
					V: []interface{}{23},
				},
			},
		},
		L:  10,
		GB: "t.name",
		OB: "t.name DESC, ot.name ASC",
		H: [][]bqb.Expr{
			{
				{
					F: "COUNT(t.name) > ?",
					V: []interface{}{1},
				},
			},
		},
	}

	q.O = 2

	sql, params, _ := q.ToPsql()

	println(sql)
	println(fmt.Sprintf("%v %T", params, params))

	q2 := &bqb.Query{}
	q2 = q2.
		Select("t.name, t.id").
		From("my_table t").
		Join("my_other_table ot ON t.id = ot.id").
		Join("users u ON t.id = u.id").
		Where(
			bqb.Valf("ST_Distance(t.geom, ot.geom) < ?", 1000),
			bqb.Valf("t.name LIKE ?", "william %"),
		).
		Where(
			bqb.Valf("ST_Distance(t.geom, GeomFromEWKT(?)) < ?", "SRID=4326;POINT(44 -111)", 1000),
		).
		Limit(10).
		Offset(2).
		OrderBy("t.name ASC, ot.name DESC").
		GroupBy("t.name").
		Having(bqb.Valf("COUNT(t.name) > ?", 2)).
		Having(bqb.Valf("COUNT(ot.name) > ?", 5))

	sql, params, _ = q2.ToPsql()

	println(sql)
	println(fmt.Sprintf("%v %T", params, params))

}
