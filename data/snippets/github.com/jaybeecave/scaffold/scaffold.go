package main

import (
	"html/template"
	"net/http"

	"bytes"

	"bufio"

	"io/ioutil"

	"os"

	"github.com/jaybeecave/render"
	errors "github.com/kataras/go-errors"
	_ "github.com/mattes/migrate/driver/postgres" //for migrations
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/urfave/cli"

	"strings"

	"github.com/serenize/snaker"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

type description struct {
	Name        string
	Method      string
	URL         string
	Description string
	Function    http.HandlerFunc
}

type descriptions []description

func (slice descriptions) Len() int {
	return len(slice)
}

func (slice descriptions) Less(i int, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice descriptions) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Field struct {
	FieldName string
	FieldType string
}

type Fields []Field

func createTable(c *cli.Context, r *render.Render, db *runner.DB) error {
	// setup
	bucket := newViewBucket()
	args := c.Args()

	if !args.Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}

	// add variables for template
	bucket.addFieldDataFromContext(c)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./models/migrations", "create_"+bucket.getStr("TableName"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "create-table", file.UpFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "drop-table", file.DownFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func addFields(c *cli.Context, r *render.Render, db *runner.DB) error {
	// setup
	bucket := newViewBucket()
	if !c.Args().Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}

	// add variables for template
	bucket.addFieldDataFromContext(c)

	file, err := migrate.Create(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./models/migrations", "fields_"+bucket.getStr("TableName"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "add-fields", file.UpFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = migrationFromTemplate(r, "remove-fields", file.DownFile, bucket)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func doMigration(c *cli.Context, r *render.Render, db *runner.DB) error {
	errs, ok := migrate.UpSync(os.Getenv("DATABASE_URL")+"?sslmode=disable", "./models/migrations")
	finalError := ""
	if !ok {
		for _, err := range errs {
			finalError += err.Error() + "\n"
		}
		return errors.New(finalError)
	}
	return nil
}

func createModel(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-model", "./models/")
}

func createRest(c *cli.Context, r *render.Render, db *runner.DB) error {
	return createSomething(c, r, db, "create-rest", "./")
}

func createSomething(c *cli.Context, r *render.Render, db *runner.DB, tmpl string, path string) error {
	bucket := newViewBucket()
	args := c.Args()

	if !args.Present() {
		// no args
		return cli.NewExitError("ERROR: No tablename defined", 1)
	}
	// add variables for template
	bucket.addFieldDataFromContext(c)

	// populate variables
	tableName := bucket.getStr("TableName")
	tableNameTitle := snaker.SnakeToCamel(tableName)
	letters := strings.Split(tableNameTitle, "")
	letters[0] = strings.ToLower(letters[0])
	tableNameCamel := strings.Join(letters, "")
	tableID := tableName + "_id"

	tnJnt := strings.Join(strings.Split(tableName, "_"), " ")

	bucket.add("TableNameSpaces", tnJnt)
	bucket.add("TableNameTitle", tableNameTitle)
	bucket.add("TableNameCamel", tableNameCamel)
	bucket.add("TableID", tableID)

	// populate more variables from column names
	columns := []*ColumnInfo{}
	err := db.Select("column_name, data_type, is_nullable").
		From("information_schema.columns").
		Where("table_schema = $1 and table_name = $2", "public", tableName).
		QueryStructs(&columns)
	if err != nil {
		return cli.NewExitError("error 10: "+err.Error(), 1)
	}
	colsDBConcat := `"`
	colsRecordPrefixedConcat := ""
	for i, col := range columns {
		if col.ColumnName == tableID {
			// we never want to include the table_id where these values are used because id's get generated from the database
			continue
		}
		colsDBConcat += col.ColumnName + `"`
		colsRecordPrefixedConcat += "record." + col.ColumnNameTitle()
		if i != (len(columns) - 1) {
			colsDBConcat += `, "`
			colsRecordPrefixedConcat += ", "
		}
	}
	bucket.add("Columns", columns)
	bucket.add("ColumnsDBStrings", template.HTML(colsDBConcat))
	bucket.add("ColumnsRecordPrefixedStrings", colsRecordPrefixedConcat)

	//
	fullpath := path + tableNameCamel + ".go.tmp"

	fo, _ := os.Create(fullpath)

	template := r.TemplateLookup(tmpl)
	wr := bufio.NewWriter(fo)
	err = template.Execute(wr, bucket.Data)
	if err != nil {
		return err
	}
	wr.Flush()
	// err = ioutil.WriteFile("./models/migrations/"+tableName+".go", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return cli.NewExitError("error 20: "+err.Error(), 1)
	}

	if err := fo.Close(); err != nil {
		return cli.NewExitError("error 30: "+err.Error(), 1)
	}
	return nil
}

type ColumnInfo struct {
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
	IsNullable string `db:"is_nullable"`
}

func (colInfo *ColumnInfo) IsNullField() bool {
	return colInfo.IsNullable == "YES"
}

func (colInfo *ColumnInfo) ColumnNameTitle() string {
	return snaker.SnakeToCamel(colInfo.ColumnName)
}

func (colInfo *ColumnInfo) ColumnType() string {
	if colInfo.DataType == "text" {
		return "string"
	}
	if colInfo.DataType == "integer" || colInfo.DataType == "numeric" {
		return "int"
	}
	if colInfo.DataType == "timestamp with time zone" {
		return "time.Time"
	}
	return ""
}

func migrationFromTemplate(r *render.Render, templateName string, file *file.File, data *viewBucket) error {
	template := r.TemplateLookup(templateName)
	buffer := bytes.NewBuffer(file.Content)
	wr := bufio.NewWriter(buffer)
	err := template.Execute(wr, data)
	if err != nil {
		return err
	}
	wr.Flush()
	err = ioutil.WriteFile(file.Path+"/"+file.FileName, buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
