package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/gophergala/food-plan-organizer/_third_party/github.com/mattn/go-sqlite3"
	"github.com/gophergala/food-plan-organizer/_third_party/github.com/rubenv/sql-migrate"
	"github.com/gophergala/food-plan-organizer/cmd/serve/manage"
	"github.com/gophergala/food-plan-organizer/cmd/serve/search"
	"github.com/gophergala/food-plan-organizer/cmd/serve/show"
)

//go:generate go-bindata -pkg main -o bindata.go migrations

var (
	listen           = flag.String("listen", ":8080", "Port to listen on")
	sr27DatabaseName = flag.String("sr27.database", "sr27.db", "SR27 Sqlite database name")
	userDatabaseName = flag.String("user.database", "user.db", "Recipe database name")
	sr27Database     *sql.DB
	userDatabase     *sql.DB
)

func newServer() http.Handler {
	var r = http.NewServeMux()
	r.Handle("/search/food/", search.NewFoodSearchServer(sr27Database))
	r.Handle("/food/", show.NewFoodShowServer(sr27Database))
	r.Handle("/recipes/", manage.NewRecipeServer(userDatabase, sr27Database))
	r.Handle("/nutrients/", show.NewNutrientShowServer(sr27Database))
	return logHandler(jsonHandler(r))
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}

func jsonHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	}
}

func runMigrations() {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	if n, err := migrate.Exec(userDatabase, "sqlite3", migrations, migrate.Up); err != nil {
		log.Printf("unable to migrate: %v", err)
	} else {
		log.Printf("Applied %d migrations!\n", n)
	}
}

func main() {
	flag.Parse()

	var err error
	if sr27Database, err = sql.Open("sqlite3", *sr27DatabaseName); err != nil {
		log.Fatalf("Failed to open sr27 database: %v\n", err)
	}
	defer sr27Database.Close()
	if userDatabase, err = sql.Open("sqlite3", *userDatabaseName); err != nil {
		log.Fatalf("Failed to open sr27 database: %v\n", err)
	}
	defer userDatabase.Close()
	runMigrations()

	log.Printf("Listening on %v", *listen)
	http.ListenAndServe(*listen, newServer())
}
