package bqb

// Dialect holds the Query dialect
type Dialect string

const (
	// PGSQL postgres dialect
	PGSQL Dialect = "postgres"
	// MYSQL MySQL dialect
	MYSQL Dialect = "mysql"
	// RAW dialect uses no parameter conversion
	RAW Dialect = "raw"
	// SQL generic dialect
	SQL Dialect = "sql"

	paramPh = "{{xX_PARAM_Xx}}"
)

// Embedded is a string type that is directly embedded into the query.
// Note: Like Embedder, this is not to be used for untrusted input.
type Embedded string

// Embedder embeds a value directly into a query string.
// Note: Since this is embedded and not bound,
// attention must be paid to sanitizing this input.
type Embedder interface {
	RawValue() string
}

// JsonMap is a custom type which tells bqb to convert the parameter to
// a JSON object without requiring reflection.
type JsonMap map[string]interface{}

// JsonList is a type that tells bqb to convert the parameter to a JSON
// list without requiring reflection.
type JsonList []interface{}
