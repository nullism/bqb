package main

import (
	"github.com/nullism/bqb"
)

func basic() {
	println("\n===[ Basic Query ]===")
	q := bqb.InsertPsql().
		Into("my_table").
		Cols("name", "age", "current_time").
		Vals(bqb.V("?, ?, ?", "someone", 42, "2021-01-01 01:01:01Z"))
	q.Print()
}

func subQuery() {
	println("\n===[ Imbedded Query ]===")
	q := bqb.InsertPsql().
		Into("my_table").
		Cols("name", "age", "current_time").
		Select(
			bqb.QueryPsql().Select("b_name", "b_age", "b_time").
				From("b_table").
				Where(bqb.V("my_age > ?", 20)).
				Limit(10),
		)
	q.Print()
}

func union() {
	println("\n===[ Union Select Query ]===")
	q := bqb.InsertPsql().
		Into("my_table").
		Cols("name", "age", "current_time").
		Union("ALL").
		Select(
			bqb.QueryPsql().Select("b_name", "b_age", "b_time").
				From("b_table").
				Where(bqb.V("my_age > ?", 20)).
				Limit(10),
		)

	q.Print()
}

func main() {
	basic()
	union()
	subQuery()
}
