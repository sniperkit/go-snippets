package resolvers

var PersonType = "Person"

var PersonSchema = `
type Person {
	name: String!
	age: Int!
}

// TODO: add fields for filter 
// input PersonFilter {
// }

` + PersonResult

var PersonResult = `
type PersonResult implements QueryResult {
	totalCount: Int!
	pageInfo: PageInfo!
	items: [Person]!
}
`

var PersonQuery = `
person(id: ID!): Person!
`

var PersonMutations = `
createPerson(id: ID!): Person!
`