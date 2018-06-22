// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/DevMine/fluxio/config"
)

const version = "1.0.0"

var (
	configPath  = flag.String("c", "", "configuration file")
	fromTable   = flag.String("r", "", "read JSON from table [name] and output to stdout")
	toTable     = flag.String("w", "", "write JSON from stdin to table [name]")
	schema      = flag.String("s", "public", "specify database schema")
	colKey      = flag.String("col-key", "key", "specify the column which acts as primary key")
	colContent  = flag.String("col-content", "content", "specify the column for json content")
	key         = flag.String("k", "", "key name for the JSON content")
	versionflag = flag.Bool("version", false, "print version.")
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	flag.Usage = func() {
		fmt.Printf("usage: %s [OPTION(S)] -c [configuration-file] [-w table-name] [-r table-name]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	if *versionflag {
		fmt.Printf("%s - %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	if len(*key) == 0 {
		err = errors.New("no key provided")
		return
	}

	if len(*configPath) == 0 {
		err = errors.New("no configuration file specified")
		return
	}

	var cfg *config.Config
	cfg, err = config.ReadConfig(*configPath)
	if err != nil {
		return
	}

	var db *sql.DB
	db, err = openDBSession(*cfg.Database)
	if err != nil {
		return
	}
	defer db.Close()

	if len(*toTable) > 0 {
		if err = verifyTable(db, *toTable, *schema, *colKey, *colContent); err != nil {
			return
		}
		if err = jsonToDB(db, *toTable, *key, *colKey, *colContent); err != nil {
			return
		}
	}

	if len(*fromTable) > 0 {
		if err = verifyTable(db, *fromTable, *schema, *colKey, *colContent); err != nil {
			return
		}
		var out string
		out, err = dbToJSON(db, *fromTable, *key, *colKey, *colContent)
		if err != nil {
			return
		}
		fmt.Println(out)
	}
}

func jsonToDB(db *sql.DB, tableName, key, colKey, colContent string) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, bufio.NewReader(os.Stdin)); err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s(%s, %s) VALUES ($1, $2)", tableName, colKey, colContent)
	if _, err := db.Exec(query, key, buf.String()); err != nil {
		return err
	}

	return nil
}

func dbToJSON(db *sql.DB, tableName, key, colKey, colContent string) (string, error) {
	query := fmt.Sprintf(`SELECT %s
                          FROM %s
                          WHERE %s=$1`, colContent, tableName, colKey)
	var out string
	if err := db.QueryRow(query, key).Scan(&out); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("cannot find content for key: '" + key + "'")
		}
		return "", err
	}

	return out, nil
}

// verifyTable makes sure that tableName from schemaName exists, has one column
// colKeyName of type string which acts as primary key and another column
// colContentName of type jsonb.
func verifyTable(db *sql.DB, tableName, schemaName, colKeyName, colContentName string) error {
	query := `SELECT EXISTS (
                  SELECT 1
                  FROM information_schema.tables
                  WHERE table_schema = $1 
                  AND table_name = $2 
              )`
	var exist bool
	if err := db.QueryRow(query, schemaName, tableName).Scan(&exist); err != nil {
		return err
	}
	if !exist {
		return errors.New("cannot find table '" + tableName + "' in schema '" + schemaName + "'")
	}

	query = `SELECT column_name
             FROM information_schema.columns
             WHERE table_schema = $1
             AND table_name = $2`
	rows, err := db.Query(query, schemaName, tableName)
	if err != nil {
		return err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		if name != colKeyName && name != colContentName {
			return errors.New("column does not exist: " + name)
		}
	}

	query = `SELECT data_type, is_nullable
             FROM information_schema.columns
             WHERE table_name = $1
             AND column_name = $2`
	var colType, nullable string
	if err := db.QueryRow(query, tableName, colKeyName).Scan(&colType, &nullable); err != nil {
		return err
	}
	if colType != "character varying" {
		return errors.New("column '" + colKeyName + "' is not of type string (" + colType + ")")
	}
	if nullable != "NO" {
		return errors.New("column '" + colKeyName + "' is nullable")
	}

	query = `SELECT data_type
             FROM information_schema.columns
             WHERE table_name = $1 
             AND column_name = $2`
	if err := db.QueryRow(query, tableName, colContentName).Scan(&colType); err != nil {
		return err
	}
	if colType != "jsonb" {
		return errors.New("column '" + colContentName + "' is not of type jsonb (" + colType + ")")
	}

	return nil
}

func openDBSession(cfg config.DatabaseConfig) (*sql.DB, error) {
	dbURL := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port=%d dbname='%s' sslmode='%s'",
		cfg.UserName, cfg.Password, cfg.HostName, cfg.Port, cfg.DBName, cfg.SSLMode)

	return sql.Open("postgres", dbURL)
}
