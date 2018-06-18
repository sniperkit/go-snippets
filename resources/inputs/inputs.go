package inputs

//Query is just the base interface
type Query interface {
	Query(key string) Value
}

//Param is just the base interface
type Param interface {
	Param(key string) Value
}

//BindBody will bind the input (POST, PUT, PATCH) body to the destination
type BindBody interface {
	BindBody(dest interface{}) error
}

//Get is for Resource GET
type Get interface {
	Query
	Param
}

//Post is for Resource POST
type Post interface {
	Query
	Param
	BindBody
}

//Patch is for Resource Patch
type Patch interface {
	Query
	Param
	BindBody
}

//Put is for Resource PUT
type Put interface {
	Query
	Param
	BindBody
}

//Delete is for Resource DELETE
type Delete interface {
	Query
	Param
}
