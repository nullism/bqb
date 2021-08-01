package main

import (
	"github.com/nullism/bqb"
)

func basic() {
	println("\n===[ Basic Query ]===")
	q := bqb.QueryPsql().
		Select("id, name, email").
		From("users").
		Where("email LIKE '%@yahoo.com'")
	q.Print()
}

func complexQuery() {
	println("\n===[ Complex Query ]===")
	q := bqb.QueryPsql().
		Select(
			"t.name", "t.id",
			bqb.V("(SELECT * FROM my_t WHERE id=?) as name", 123),
		).
		From("my_table t").
		Join("my_other_table ot ON t.id = ot.id").
		Join(bqb.V("users u ON t.id > ?", 7)).
		Where(
			bqb.Or(
				bqb.V("ST_Distance(t.geom, ot.geom) < ?", 100),
				bqb.V("t.name LIKE ?", "william%"),
				bqb.V("ST_Distance(t.geom, GeomFromEWKT(?)) < ?", "SRID=4326;POINT(44 -111)", 100),
			),
		).
		OrderBy("t.name ASC, ot.name DESC").
		GroupBy("t.name", "t.id").
		Having(
			bqb.V("COUNT(t.name) > ?", 1),
			bqb.V("COUNT(ot.name) > ?", 1),
		).
		Offset(15).
		Limit(10)

	q.Print()
}

func join() {
	println("\n===[ Join Query ]===")
	q := bqb.QueryPsql().
		Select("uuidv3_generate() as uuid", "u.id", "UPPER(u.name) as screamname", "u.age", "e.email").
		From("users u").
		Join("emails e ON e.user_id = u.id").
		Where(
			bqb.Or(
				bqb.And(
					bqb.V("u.id IN (?, ?, ?)", 1, 3, 5),
					bqb.V("AND e.email LIKE ?", "%@gmail.com"),
				),
				bqb.And(
					bqb.V("u.id IN (?, ?, ?)", 2, 4, 6),
					bqb.V("AND e.email LIKE ?", "%@yahoo.com"),
				),
				bqb.V("u.id IN (?)", []int{7, 8, 9, 10, 11, 12}),
			),
		).
		OrderBy("u.age DESC").
		Limit(10)
	q.Print()
}

func andOr() {
	println("\n===[ And/Or Query ]===")
	q := bqb.QueryPsql().Select("*").From("patrons").
		Where(
			bqb.Or(
				bqb.And(
					"drivers_license IS NOT NULL",
					bqb.And("age > 20", "age < 60"),
				),
				bqb.And(
					"drivers_license IS NULL",
					"age >= 60",
				),
				"is_known = true",
			),
		)
	q.Print()
}

func raw() {
	println("\n===[ Raw Dialect Query ]===")
	q := bqb.QueryRaw().
		Select("*", bqb.V("my_function(?, ?)", "name", true)).
		From("my_table").
		Where(
			bqb.V("my_value = ?", 1234),
			bqb.V("my_other_value IS ?", nil),
		)
	q.Print()
}

func v() {
	println("\n===[ Bind Query ]===")
	email := "foo@bar.com"
	password := "p4ssw0rd"
	q := bqb.QueryPsql().
		Select("*").
		From("users").
		Where(
			bqb.And(
				bqb.V("email = ?", email),
				bqb.V("password = ?", password),
			),
		)
	q.Print()
}

func main() {

	basic()
	v()
	join()
	complexQuery()
	andOr()
	raw()
}
