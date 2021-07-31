package bqb

const (
	pgsql = "postgres"
	mysql = "mysql"
)

type Expr struct {
	F string
	V []interface{}
}

type Query struct {
	SE []Expr
	F  string
	FE []Expr
	J  []string
	W  [][]Expr
	OB string
	L  int
	O  int
	GB string
	H  [][]Expr
}
