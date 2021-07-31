package main

import (
	"github.com/nullism/bqb"
)

func basic() {
	q := bqb.New(bqb.PGSQL).
		Select("id, name, email").
		From("users").
		Where("email LIKE '%@yahoo.com'")
	q.Print()
}

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
	q.Print()
}

func complexQuery() {
	q := bqb.New(bqb.PGSQL).
		Select("t.name", "t.id", bqb.Valf("(SELECT * FROM my_t WHERE id=?) as name", 123)).
		From("my_table t").
		Join("my_other_table ot ON t.id = ot.id").
		Join(bqb.Valf("users u ON t.id = ?", 7)).
		Where(
			bqb.Valf("ST_Distance(t.geom, ot.geom) < ?", 101),
			bqb.Valf("t.name LIKE ?", "william%"),
			bqb.Valf("ST_Distance(t.geom, GeomFromEWKT(?)) < ?", "SRID=4326;POINT(44 -111)", 102),
		).
		Limit(10).
		Offset(2).
		OrderBy("t.name ASC, ot.name DESC").
		GroupBy("t.name", "t.id").
		Having(
			bqb.Valf("COUNT(t.name) > ?", 2),
			bqb.Valf("COUNT(ot.name) > ?", 5),
		)

	q.Print()
}

func andOr() {
	q := bqb.New(bqb.PGSQL).Select("*").From("patrons").
		Where(
			"(drivers_license IS NOT NULL AND (age > 20 AND age < 60))",
			"(drivers_license IS NULL AND age > 60)",
			"is_known = true",
		)
	q.Print()
}

func valf() {
	email := "foo@bar.com"
	password := "p4ssw0rd"
	q := bqb.New(bqb.PGSQL).
		Select("*").
		From("users").
		Where(
			bqb.Valf("email = ?", email),
			bqb.Valf("password = ?", password),
		)
	q.Print()
}

func main() {

	// rawApi()
	// complexQuery()
	andOr()
	// basic()
	// valf()

	q := bqb.New(bqb.PGSQL).
		Select("*").From("users").Where(
		bqb.Valf(
			"(name LIKE ? AND age < ?) OR (name LIKE ? AND (age > ? AND email LIKE ?)",
			"smith%", 20, "barb%", 22, "%@gmail.com",
		),
	)
	q.Print()

}
