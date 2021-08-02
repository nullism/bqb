package main

import (
	"github.com/nullism/bqb"
)

func basic() {
	println("\n===[ Basic Update ]===")
	q := bqb.Update("my_table").
		Set(
			bqb.V("name = ?", "McCallister"),
			"age = 20", "current_time = CURRENT_TIMESTAMP()",
		).
		Where(
			bqb.V("name = ?", "Mcallister"),
		).Postgres()
	q.Print()
}

func subquery() {
	println("\n===[ Advanced Update ]===")

	timeQ := bqb.QueryPsql().Select("timestamp").
		From("time_data").Where(bqb.V("is_current = ?", true)).
		Limit(1)

	nameQ := bqb.QueryPsql().
		Select("name").
		From("users").
		Where(bqb.V("name LIKE ?", "%allister"))

	q := bqb.Update("my_table").
		Set(
			bqb.V("name = ?", "McCallister"),
			"age = 20",
			bqb.Concat(
				"current_timestamp = ",
				timeQ.Enclose(),
			),
		).
		Where(
			bqb.Concat(
				"name IN ",
				nameQ.Enclose(),
			),
		).Postgres()
	q.Print()
}

func main() {
	basic()
	subquery()
}
