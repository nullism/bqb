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

func selects() {
	println("\n===[ Basic Query ]===")
	q := bqb.InsertPsql().
		Into("my_table").
		Cols("name", "age", "current_time").
		Select("other_name", "other_age", "other_time").
		From("other_table").
		Where(bqb.V("other_age IS NOT ?", nil))
	q.Print()
}

func main() {
	basic()
	selects()
}
