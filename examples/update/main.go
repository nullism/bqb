package main

import (
	"github.com/nullism/bqb"
)

func basic() {
	println("\n===[ Basic Update ]===")
	q := bqb.UpdatePsql().
		Update("my_table").
		Set(
			bqb.V("name = ?", "McCallister"),
			"age = 20", "current_time = CURRENT_TIMESTAMP()",
		).
		Where(
			bqb.V("name = ?", "Mcallister"),
		)
	q.Print()
}

func subquery() {
	println("\n===[ Advanced Update ]===")

	timeQ := bqb.QueryPsql().Select("timestamp").
		From("time_data").Where("is_current = true").
		Limit(1)

	nameQ := bqb.QueryPsql().
		Select("name").
		From("users").
		Where(bqb.V("name LIKE ?", "%allister"))

	q := bqb.UpdatePsql().
		Update("my_table").
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
		)
	q.Print()
}

func main() {
	basic()
	subquery()
}
