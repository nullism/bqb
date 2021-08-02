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
	q := bqb.UpdatePsql().
		Update("my_table").
		Set(
			bqb.V("name = ?", "McCallister"),
			"age = 20", "current_time = CURRENT_TIMESTAMP()",
		).
		Where(
			"name IN (",
			bqb.QueryPsql().Select("name").From("users").Where("name LIKE '%allister'"),
			")",
		)
	q.Print()
}

func main() {
	basic()
	subquery()
}
