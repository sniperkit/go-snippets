package responses

//Response is the generic response that Resource methods (Get, Post, etc) will return
type Response struct {
	Status Status
	Body   interface{}
	Format Format
}
