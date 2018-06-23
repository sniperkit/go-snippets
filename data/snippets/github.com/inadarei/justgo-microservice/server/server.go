package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/**
 * StartServer is the main entrance into the server.
 */
func StartServer(serverPort string) {
	go startHealthServerNative(serverPort)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/", homeEndpoint)
	router.Run(":" + serverPort)

}

func homeEndpoint(c *gin.Context) {
	toGreetName := "Go enthusiasts"
	c.String(http.StatusOK, "Hello %s! This is pretty awesome! Hot-reloading is great!", toGreetName)
}

//---- End of the Gin-based implementation of the home endpoint
//----
// Note: the native (non Gin-based) implementation below is hidden
// from docker-compose and the code is only shown here as a demonstration of
// how to implement things without Gin, if you wanted to do so
func startHealthServerNative(serverPort string) {
	http.HandleFunc("/health", viewHandlerHealth)
	healthPort, errParse := strconv.Atoi(serverPort)
	healthPort = healthPort + 1
	healthPortStr := strconv.Itoa(healthPort)
	abortIfErr(errParse)
	err := http.ListenAndServe(":"+healthPortStr, nil)
	if err != nil {
		log.Fatal("ERROR: couldn't start server: ", err)
	} else {
		log.Info("Healthcheck started successfully at: ", healthPortStr)
	}
}

func viewHandlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/health+json")
	response := make(map[string]string)
	response["status"] = "pass"

	res, err := json.MarshalIndent(response, "", "  ")
	abortIfErr(err)
	fmt.Fprintf(w, string(res))
}

// Simple exit if error, to avoid putting same 4 lines of code in too many places
func abortIfErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
