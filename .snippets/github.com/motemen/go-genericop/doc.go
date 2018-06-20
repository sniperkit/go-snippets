// Package genericop provides functions that does operators (e.g. "+" or ">")
// on interface{} args. Each function returns the result along with an error,
// which indicates if the argument types were incompatible on the operation.
// You can use MustT(x, err) family to assert the resulting type.
package genericop
