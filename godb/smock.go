package godb

import (
	"database/sql"
	"log"
)

type user struct {
	id   int
	name string
}

func getUser(db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM user")
	if err != nil {
		log.Printf("failed to fetch user, err: %v ", err)
		return
	}
	defer rows.Close()

	var users []*user
	for rows.Next() {
		s := &user{}
		if err := rows.Scan(&s.id, &s.name); err != nil {
			log.Printf("failed to scan user, err: %v", err)
			return
		}
		users = append(users, s)
	}

	if rows.Err() != nil {
		log.Printf("failed to read all posts, err: %v", rows.Err())
		return
	}

	for k, v := range users {
		log.Printf("%d ==> %v\n", k, v)
	}
}
