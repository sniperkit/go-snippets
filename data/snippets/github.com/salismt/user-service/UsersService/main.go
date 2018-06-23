package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	CreateUsersQueue = "CREATE_USER"
	UpdateUsersQueue = "UPDATE_USER"
	DeleteUsersQueue = "DELETE_USER"
)

func main() {
	config := BaseConfig{}
	config.Load()

	fmt.Printf("Port is %s", config.GetValue("port"))
	fmt.Printf("Redis is %s", config.GetValue("APP_RD_ADDRESS"))

	var numWorkers int
	cache := Cache{Enable: true}

	// flag to receive information directly from the command line
	flag.StringVar(
		&cache.Address,
		"redis_address",
		config.GetValue("APP_RD_ADDRESS"),
		"Redis Address",
	)

	flag.StringVar(
		&cache.Auth,
		"redis_auth",
		config.GetValue("APP_RD_AUTH"),
		"Redis Auth",
	)

	flag.StringVar(
		&cache.DB,
		"redis_db_name",
		config.GetValue("APP_RB_DBNAME"),
		"Redis DB Name",
	)

	flag.IntVar(
		&cache.MaxActive,
		"redis_max_active",
		60,
		"Redis Max Active",
	)

	flag.IntVar(
		&cache.IdleTimeoutSecs,
		"redis_timeout",
		60,
		"Redis timeout in seconds",
	)

	flag.IntVar(
		&numWorkers,
		"num_workers",
		10,
		"Numer of workers to consume queue",
	)

	flag.Parse()
	cache.Pool = cache.NewCachePool()

	connectionString := os.Getenv("DATABASE_DEV_URL")

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	go UsersToDB(numWorkers, db, &cache, CreateUsersQueue)
	go UsersToDB(numWorkers, db, &cache, UpdateUsersQueue)
	go UsersToDB(numWorkers, db, &cache, DeleteUsersQueue)

	a := App{}
	a.Initialize(cache, db)
	a.Run(fmt.Sprintf(":%s", config.GetValue("port")))
}
