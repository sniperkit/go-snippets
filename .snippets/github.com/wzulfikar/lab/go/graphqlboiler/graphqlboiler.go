// generate graphql resolvers and schema
// based on struct
//
// sample code:
//
// 	graphqlboiler.Boil(graphqlboiler.Tpl{
//		RootResolver: "RootResolver",
//		schema:       Person{},
//	}, "./graphqlboiler/")
package graphqlboiler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/wzulfikar/lab/go/graphqlboiler/templates"
	"gopkg.in/volatiletech/null.v6"
)

// sample schema
type SampleSchemaPerson struct {
	Name      string
	Age       int
	Hobby     null.String
	Birthdate null.Time

	married bool `json:"-"`
}

type resolverField struct {
	name      string
	fieldType string
}

type Tpl struct {
	RootResolver string
	Schema       interface{}
	Repo         string
}

func Boil(tpl Tpl, path string) {
	typeName := string(reflect.TypeOf(tpl.Schema).Name())
	resolverName := strings.ToLower(typeName[0:1]) + typeName[1:] + "Resolver"
	fields := reflectFields(tpl.Schema)

	model := strings.ToUpper(typeName[0:1]) + typeName[1:]
	modelPlural := inflection.Plural(model)

	data := struct {
		HasScalar bool

		TypeName,
		TypeNameLowerCase,
		SchemaFields,
		ModelPlural,
		ResolverName,
		RootResolver,
		Repo string
	}{
		false,
		typeName,
		strings.ToLower(typeName[0:1]) + typeName[1:],
		schemaFields(fields),
		modelPlural,
		resolverName,
		tpl.RootResolver,
		tpl.Repo,
	}

	fmt.Println("Generating graphqlboiler for struct `" + typeName + "`")

	tplResolverFields := tplResolverFields(resolverName, fields)
	fmt.Println("✔ Resolver fields template")

	data.HasScalar = strings.Contains(tplResolverFields, "scalar")

	tplResolver, err := parseTpl(templates.Resolver, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✔ Resolver template")

	tplMutations, err := parseTpl(templates.Mutations, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✔ Mutations template")

	tplQueryResult, err := parseTpl(templates.QueryResult, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✔ Query result template")

	tplSchema, err := parseTpl(templates.Schema, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✔ Schema template")

	mustWrite(path+strings.ToLower(typeName)+".go", tplResolver+tplResolverFields)
	mustWrite(path+strings.ToLower(typeName)+"_result.go", tplQueryResult)
	mustWrite(path+strings.ToLower(typeName)+"_schema.go", tplSchema)
	mustWrite(path+strings.ToLower(typeName)+"_mutations.go", tplMutations)
	fmt.Println("✔ Write templates to file")
	fmt.Println("[DONE!]")
}

func tplResolverFields(resolverName string, fields []resolverField) string {
	var resolvers string
	var fldType string

	for _, field := range fields {
		fldName := field.name
		returnObj := `r.o.` + fldName
		preReturnCode := ""

		switch field.fieldType {
		case "Uint", "uint", "Int", "int":
			fldType = "int32"
			returnObj = `int32(r.o.` + fldName + `)`

		case "null.Uint":
			fldType = "*int32"
			preReturnCode = `if !r.o.` + fldName + `.Valid {
		return nil, nil
	}
	v := int32(r.o.` + fldName + `.Uint)
	`
			returnObj = `&v`

		case "null.Int":
			fldType = "*int32"
			preReturnCode = `if !r.o.` + fldName + `.Valid {
		return nil, nil
	}
	v := int32(r.o.` + fldName + `.Int)
	`
			returnObj = `&v`

		case "decimal":
			fldType = "scalar.Decimal"

		case "time", "Time":
			fldType = "scalar.Time"
			returnObj = `scalar.Time{r.o.` + fldName + `}`

		case "null.Time":
			fldType = "*scalar.Time"
			preReturnCode = `if !r.o.` + fldName + `.Valid {
		return &scalar.Time{}, nil
	}
	`
			returnObj = `&scalar.Time{r.o.` + fldName + `.Time}`

		case "null.String":
			fldType = "*string"
			returnObj = `r.o.` + fldName + `.Ptr()`

		default:
			fldType = field.fieldType
		}

		if strings.ToUpper(fldName) == "ID" {
			if fldType == "int" {
				returnObj = `graphql.ID(strconv.Itoa(r.o.` + fldName + `))`
			} else {
				returnObj = `graphql.ID(r.o.` + fldName + `)`
			}
			fldType = "graphql.ID"
		} else if strings.HasSuffix(fldName, "ID") {
			fldName = strings.Replace(fldName, "ID", "", 1)

			o := strings.ToUpper(fldName[:1]) + fldName[1:]
			resolver := strings.ToLower(fldName[:1]) + fldName[1:] + "Resolver"

			fldType = "*" + resolver
			preReturnCode = `if r.o.R == nil || r.o.R.` + o + ` == nil{
		r.o.L.Load` + o + `(r.rr.Db, true, r.o)
	}
	`
			returnObj = `&` + resolver + `{rr: r.rr, o: r.o.R.` + o + `}`
		}

		resolvers += `
func (r *` + resolverName + `) ` + fldName + `() (` + fldType + `, error) {
	` + preReturnCode + `return ` + returnObj + `, nil
}
`
	}
	return strings.TrimRight(resolvers, "\n") + "\n"
}

func reflectFields(schema interface{}) (fields []resolverField) {
	rv := reflect.ValueOf(schema)
	for i := 0; i < rv.NumField(); i++ {
		fld := rv.Field(i)
		typeFld := rv.Type().Field(i)

		// don't reflect fields with json tag "-"
		if typeFld.Tag.Get("json") == "-" {
			continue
		}

		fldType := fld.Type().Name()
		if fld.Kind() == reflect.Struct {
			if fld.FieldByName("Valid").IsValid() {
				fldType = "null." + fldType
			}
		}

		fields = append(fields, resolverField{typeFld.Name, fldType})
	}
	return fields
}

func schemaFields(fields []resolverField) string {
	var schemaFields string
	for _, field := range fields {
		fldType := strings.ToUpper(field.fieldType[0:1]) + field.fieldType[1:]
		fldName := strings.ToLower(field.name[0:1]) + field.name[1:]

		if strings.HasPrefix(fldType, "Uint") {
			fldType = strings.Replace(fldType, "Uint", "Int", 1)
		}

		// determine if a field is required or not
		if strings.HasPrefix(fldType, "Null.") {
			fldType = strings.Replace(fldType, "Null.", "", 1)
		} else {
			fldType += "!"
		}

		// field `id` should return graphql `ID!`
		if strings.ToUpper(fldName) == "ID" {
			fldName = "id"
			fldType = "ID!"
		} else if strings.HasSuffix(fldName, "ID") {
			fldName = strings.Replace(fldName, "ID", "", 1)
			fldType = strings.ToUpper(fldName[:1]) + fldName[1:]
		}

		schemaFields += fmt.Sprintf("\t%s: %s\n", fldName, fldType)
	}

	return strings.TrimRight(strings.TrimLeft(schemaFields, "\t"), "\n")
}

func parseTpl(tplString string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tplString)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func mustWrite(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot create file: %v", err))
	}
	defer file.Close()
	fmt.Fprintf(file, content)
}
