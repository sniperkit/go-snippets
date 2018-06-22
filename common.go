package main

const TagSep = "\u001F"

// Id is int32 so it matches INTEGER (SERIAL is INTEGER)
// cf. http://www.postgresql.org/docs/9.3/static/datatype-numeric.html
// Doc allows data exchange between DB and User.
// Capitalized name for JSON.
type Doc struct {
	Id			int32
	Name		string
	Type		string
	Content		string
	Uid			int32		// Id given by Auth
	Tags		[]string
}
