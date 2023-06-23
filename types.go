package bqb

// Embedder embeds a value directly into a query string.
// Note: Since this is embedded and not bound,
// attention must be paid to sanitizing this input.
type Embedder interface {
	RawValue() string
}
