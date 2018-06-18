package templates

const Schema = `package resolvers

var {{.TypeName}}Type = "{{.TypeName}}"

var {{.TypeName}}Schema = ` + "`" + `
type {{.TypeName}} {
	{{.SchemaFields}}
}

// TODO: add fields for filter 
// input {{.TypeName}}Filter {
// }

` + "` + {{.TypeName}}Result" + `

var {{.TypeName}}Result = ` + "`" + `
type {{.TypeName}}Result implements QueryResult {
	totalCount: Int!
	pageInfo: PageInfo!
	items: [{{.TypeName}}]!
}
` + "`" + `

var {{.TypeName}}Query = ` + "`" + `
{{.TypeNameLowerCase}}(id: ID!): {{.TypeName}}!
` + "`" + `

var {{.TypeName}}Mutations = ` + "`" + `
create{{.TypeName}}(id: ID!): {{.TypeName}}!
` + "`"
