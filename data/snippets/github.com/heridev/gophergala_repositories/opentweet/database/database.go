package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophergala/opentweet/protocol"
	"log"
	"time"
)

type DB struct {
	db *sql.DB
}

func NewDB() (DB, error) {
	var newDB DB
	db, err := sql.Open("mysql", "root:mysecretpassword@tcp(mysql:3306)/")
	//db, err := sql.Open("mysql", "root:mysecretpassword@tcp(localhost:3306)/")
	if err != nil {
		return newDB, fmt.Errorf("Error opening db: %v", err)
	}

	_, err = db.Exec(
		"CREATE DATABASE IF NOT EXISTS opentweet;")
	if err != nil {
		return newDB, fmt.Errorf("Error creating db: %v", err)
	}

	_, err = db.Exec(
		"USE opentweet;")
	if err != nil {
		return newDB, fmt.Errorf("Error using db: %v", err)
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS users " +
			"(name VARCHAR(32) PRIMARY KEY, password VARCHAR(32));")
	if err != nil {
		return newDB, fmt.Errorf("Error creating table: %v", err)
	}

	log.Printf("Database connected and setup")
	newDB.db = db
	return newDB, nil
}

func (db DB) GetTweets(name string, from, to time.Time) ([]protocol.Tweet, error) {
	//sql injection here
	query := fmt.Sprintf("SELECT time, tweet FROM `%v` WHERE time BETWEEN ? AND ?;", name)
	rows, err := db.db.Query(
		query,
		from.Unix(),
		to.Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error getting users tweets: %v", err)
	}

	tweets := make([]protocol.Tweet, 0)
	for rows.Next() {
		var timeStamp int64
		var tweet string
		err = rows.Scan(&timeStamp, &tweet)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get tweets from db: %v", err)
		}
		tweets = append(tweets, protocol.Tweet{time.Unix(timeStamp, 0), tweet})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row error: %v", err)
	}
	return tweets, nil
}

func (db DB) RegisterUser(name, password string) error {
	log.Printf("Registering new user: %v", name)

	_, err := db.db.Exec(
		"INSERT INTO users (name, password)" +
			"VALUES (?, ?);",
		name,
		password,
	)
	if err != nil {
		return fmt.Errorf("Error inserting user: %v", err)
	}

	//sql injection here
	query := fmt.Sprintf("CREATE TABLE `%v` (time BIGINT PRIMARY KEY, tweet VARCHAR(256))", name)
	
	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error creating user tweet table: %v", err)
	}

	return nil
}

func (db DB) PostTweet(user, password, tweet string) error {
	rows, err := db.db.Query(
		"SELECT password FROM users WHERE name = ?;", user)
	if err != nil {
		return fmt.Errorf("Error finding user: %v", err)
	}
	//should be 1 row
	ok := rows.Next()
	if ok != true {
		err = rows.Err()
		return fmt.Errorf("Error geting user row: %v", err)
	}
	var dbPass string
	err = rows.Scan(&dbPass)
	if err != nil {
		return fmt.Errorf("Error getting user password: %v", err)
	}
	rows.Close()

	if password != dbPass {
		return fmt.Errorf("Error invalid password")
	}

	//sql injection here
	query := fmt.Sprintf("INSERT INTO `%v` (time, tweet) VALUES (?, ?);", user)
	_, err = db.db.Exec(
		query,
		time.Now().Unix(),
		tweet,
	)
	if err != nil {
		return fmt.Errorf("Error inserting tweet: %v", err)
	}

	log.Printf("user %v posted tweet %v", user, tweet)
	return nil
}
