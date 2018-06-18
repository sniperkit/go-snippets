package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type State struct {
	ID   int
	Name string
}

type Park struct {
	ID          int
	Name        string
	Description string
	NearestCity string
	Visitors    int
	Established time.Time
	StateID     int
	State       State
}

func main() {

	const (
		host     = "172.17.0.2"
		port     = 5432
		user     = "postgres"
		password = "P@ssw0rd!"
		dbname   = "national_parks"
	)

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("DB Conn error: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping error: ", err)
	}
	defer db.Close()
	query := `
    Select 
        p.id, 
        p.name, 
        p.description, 
        p.nearest_city, 
        p.visitors, 
        p.established, 
        p.state_id, 
        s.id, 
        s.name
    FROM parks as p 
        LEFT OUTER JOIN states as s 
        ON s.ID=p.state_id;`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("DB Query error: ", err)
	}
	defer rows.Close()

	var parks []*Park
	for rows.Next() {
		var p = &Park{}
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.NearestCity,
			&p.Visitors,
			&p.Established,
			&p.StateID,
			&p.State.ID,
			&p.State.Name); err != nil {
			log.Fatal(err)
		}
		parks = append(parks, p)
	}

	for _, p := range parks {
		fmt.Printf("%v", p)
	}
}
