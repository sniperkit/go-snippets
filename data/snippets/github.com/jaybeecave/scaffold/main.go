package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	dat "gopkg.in/mgutz/dat.v1"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"

	"github.com/jaybeecave/render"
	dotenv "github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	dotenv.Load() // load from .env file where scaffold is run
	render := getRenderer()
	db := getDBConnection()
	app := cli.NewApp()
	app.Name = "scaffold"
	app.Usage = "generate models & migrations using dat"
	app.Commands = []cli.Command{
		{
			Name:    "table",
			Aliases: []string{"t"},
			Usage:   "Create a new table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return createTable(c, render, db)
			},
		},
		{
			Name:    "fields",
			Aliases: []string{"f"},
			Usage:   "Add fields to an existing table [tablename] [fieldname:fieldtype]",
			Action: func(c *cli.Context) error {
				return addFields(c, render, db)
			},
		},
		{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "create a model from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createModel(c, render, db)
			},
		},
		{
			Name:    "rest",
			Aliases: []string{"r"},
			Usage:   "create a restful interface from a table [tablename]",
			Action: func(c *cli.Context) error {
				return createRest(c, render, db)
			},
		},
		{
			Name:    "migration",
			Aliases: []string{"mi"},
			Usage:   "perform schema migration",
			Action: func(c *cli.Context) error {
				return doMigration(c, render, db)
			},
		},
	}
	app.Run(os.Args)
}

func getRenderer() *render.Render {
	r := render.New(render.Options{
		Directory: "./models/templates",
	})
	return r
}

func getDBConnection() *runner.DB {
	//get url from ENV in the following format postgres://user:pass@192.168.8.8:5432/spaceio")
	dbURL := os.Getenv("DATABASE_URL")
	u, err := url.Parse(dbURL)
	if err != nil {
		panic(err)
	}

	username := u.User.Username()
	pass, isPassSet := u.User.Password()
	if !isPassSet {
		panic("no database password")
	}
	host, port, _ := net.SplitHostPort(u.Host)
	dbName := strings.Replace(u.Path, "/", "", 1)

	db, _ := sql.Open("postgres", "dbname="+dbName+" user="+username+" password="+pass+" host="+host+" port="+port+" sslmode=disable")
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	// db connection
	return runner.NewDB(db, "postgres")
}

// for storing variables when running the templates
type viewBucket struct {
	Data map[string]interface{}
}

func newViewBucket() *viewBucket {
	return &viewBucket{Data: map[string]interface{}{
		"LTEqStr": template.HTML(`<=`),
		"GTEqStr": template.HTML(`>=`),
		"LTStr":   template.HTML(`<`),
		"GTStr":   template.HTML(`>`),
	}}
}

func (viewBucket *viewBucket) add(key string, value interface{}) {
	viewBucket.Data[key] = value
}

func (viewBucket *viewBucket) getStrSafe(key string) (string, error) {
	val := viewBucket.Data[key]
	if val == nil {
		return "", errors.New("could not find " + key)
	}
	strVal, ok := val.(string)
	if !ok {
		return "", errors.New("could not cast " + key + " to string")
	}
	return strVal, nil
}

// getStr - returns a string for the provided key. Will panic if key not found
func (viewBucket *viewBucket) getStr(key string) string {
	val, err := viewBucket.getStrSafe(key)
	if err != nil {
		panic(err)
	}
	return val
}

func (viewBucket *viewBucket) addFieldDataFromContext(c *cli.Context) {
	args := c.Args()
	viewBucket.add("TableName", args.First())

	fields := Fields{}
	for _, arg := range args {
		if args.First() == arg {
			continue // we dont care about the first arg as its the TableName
		}
		if strings.Contains(arg, ":") {
			strSlice := strings.Split(arg, ":")
			field := Field{
				FieldName: strSlice[0],
				FieldType: strSlice[1],
			}
			fields = append(fields, field)
		}
	}
	viewBucket.add("Fields", fields)
}
