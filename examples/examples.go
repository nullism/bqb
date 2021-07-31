package main

import (
	"github.com/nullism/bqb"
)

func rawApi() {
	q := bqb.Query{
		SE: []bqb.Expr{{F: "*, t.name, t.id"}},
		FE: []bqb.Expr{{F: "my_table t"}},
		JE: []bqb.Expr{{F: "my_other_tab ot ON ot.id = t.id"}},
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
		GB: []bqb.Expr{{F: "t.name"}},
		OB: []bqb.Expr{{F: "t.name DESC, ot.name ASC"}},
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
	q.Print(bqb.PGSQL)
}

func complexQuery() {
	q := &bqb.Query{}
	q = q.
		Select("t.name", "t.id", bqb.Valf("(SELECT * FROM my_t WHERE id=?) as name", 123)).
		From("my_table t").
		Join("my_other_table ot ON t.id = ot.id").
		Join(bqb.Valf("users u ON t.id = ?", 7)).
		Where(
			bqb.Valf("ST_Distance(t.geom, ot.geom) < ?", 101),
			bqb.Valf("t.name LIKE ?", "william%"),
		).
		Where(
			bqb.Valf("ST_Distance(t.geom, GeomFromEWKT(?)) < ?", "SRID=4326;POINT(44 -111)", 102),
			"the_fox > the_hound",
		).
		Limit(10).
		Offset(2).
		OrderBy("t.name ASC, ot.name DESC").
		GroupBy("t.name").
		Having(bqb.Valf("COUNT(t.name) > ?", 2)).
		Having(bqb.Valf("COUNT(ot.name) > ?", 5))

	q.Print(bqb.PGSQL)
}

func main() {

	// rawApi()
	complexQuery()

	q1 := &bqb.Query{}
	q1.Select("t.name", "*").From("table1 t").Where(bqb.Valf("t.name LIKE ?", "william%"))
	q1.Print(bqb.RAW)

}
