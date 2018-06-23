---
title: "GraphQL + SQLBoiler Generator"
date: 2018-03-05T00:21:19+08:00
tags: ["go", "sqlboiler", "graphql"]
draft: false
---

This note as a continuation of my journey with Go template. Based on the knowledge of what template package can do, we'll generate a resolvers code from a given sqlboiler model struct to to satisfy graphql implementation of neelance/graphql-go.

<!--more-->

## Synopsis

I'm working on a project that uses MySQL as database, and wanted to have GraphQL implemented. To handle the database, we'll use SQLBoiler from [volatiletech/sqlboiler](https://github.com/volatiletech/sqlboiler) (they've very good documentation!). As for the GraphQL, we'll use [neelance/graphql-go](http://github.com/neelance/graphql-go). 

Now, the flow to implement the GraphQL is:

1. GraphQL endpoint starts with a root resolver, ie. `RootResolver`. A sample code would look like this:
    
    ```go
    // import graphql "http://github.com/neelance/graphql-go"
    
    resolver := &RootResolver{}
    var schema *graphql.Schema = graphql.MustParseSchema(resolvers.Schema, resolver)
    ```

2. Every endpoint must have resolver. Assuming we want to create a graphql endpoint for `user`, then we'll have a `User()` method attached to `RootResolver`, that will return a `userResolver`
3. Every endpoint resolver (ie. `userResolver`), will resolve every available fields. 
If a `user` has 3 fields `name`, `email` and `username`, `userResolver` must have 3 accompanying methods: `Name()`, `Email()`, and `Username()`
4. Once `userResolver` is completed, update the graphql schema.

We can represent the relationship of methods mentioned above like so:
```
▾ RootResolver
    .User()     *userResolver
▾ *userResolver
    .Name()     string
    .Email()    string
    .Username() string
```

## Generator Approach

Looking at the relationships, it will be a repetitive work to hand-code all the resolvers for each model. Thus, we'll come with a generator approach, that will take a struct, generate those resolvers code based on the struct, and store the generated codes to files inside the directory we choose. 

## The Code

```go
// import "github.com/wzulfikar/lab/go/graphqlboiler"

type User struct {
    Name      string
    Email     string
    Username  string
}

resolversDir := "/path-to-project/resolvers/"
graphqlboiler.Boil(graphqlboiler.Tpl{
    RootResolver: "RootResolver",
    Schema:       User{},
    Repo:         "gitlab.com/wzulfikar/iiumpayment",
}, resolversDir)
```

When above code is executed, the `graphqlboiler.Boil` will use Go's reflection ability from [pkg/reflect](https://golang.org/pkg/reflect/) to get the struct type, field names and field types. In this case, it will know that the struct's type name is `User`, and it has 3 fields; `Name`, `Email` and `Username`, which its types are all `string`. 

Using this information and common convention, the generator determines that the resolver name is `userResolver`, and proceeds with generating resolvers code for the `User` struct. Upon completion, there will be four new files stored in resolvers directory: 

```
▾ path-to-project/
    ▾ resolvers/
        user.go
        user_mutations.go
        user_result.go
        user_schema.go
```

Above four files contains the necessary codes to run the graphql;

1. `user.go`
    - attach `userResolver` to `RootResolver`, and 
    - attach methods for `userResolver` (`Name()`, `Email()` and `Username()`)
2. `user_mutations.go`
    - attach mutation endpoint `CreateUser()` to `RootResolver`.
3. `user_result.go`
    - Contains complimentary codes to add pagination for the query
4. `user_schema.go`
    - the GraphQL schema os `User`, with exported variables

Using the generator, we've completed steps 1 to 3 of the flow presented in [#Synopsis](#synopsis). All the repetitive steps are done by generator. We can proceed directly to step 4: updating our graphql schema, which basically just putting the variables from `user_schema.go` to the schema that will be parsed into `graphql.MustParseSchema()`.

Once the schema is updated, we can run the test and spin up our graphql server!

## Closing

The templates used to generate above codes are available at [wzulfikar/lab/graphqlboiler/templates](https://github.com/wzulfikar/lab/tree/master/go/graphqlboiler/templates). While the generator itself is by any mean not a sophisticated code, it has helped me in building GraphQL endpoint from the same scenario (a Go project with SQLBoiler and neelance/graphql-go). If you'd like to explore other graphql-go generator, you may want to see this: https://github.com/vektah/gqlgen.

Lastly, the full code for the generator is available here: https://github.com/wzulfikar/lab/tree/master/go/graphqlboiler

***Till next. See ya!***