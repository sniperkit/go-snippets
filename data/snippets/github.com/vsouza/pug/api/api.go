package api

import (
	"log"

	"github.com/codegangsta/negroni"
	"github.com/vsouza/pug/config"
)

// Resources is an exported struct to store all service resources.
// Dabatabase connections, logger handlers, metrics services and more.
// All the resources that will be used by handlers and models.
// They must be stored in this struct.
type Resources struct {
	DB      *dbResource
	Metrics *metricsResource
	Logger  *handler
}

// StartService is an exported method to start all crucial components of service,
// router, handlers, loggers, middleware, etc.
// If any error occurred, a fatal error is thrown.
func StartService(c *config.Config) {
	db, err := NewDBResource(c)
	if err != nil {
		log.Fatal("db resource error")
	}
	resources := &Resources{DB: db}
	router := NewRouter(c, resources)
	router.InitRoutes()
	n := negroni.New(negroni.NewRecovery(), NewLogger())
	n.UseHandler(router.GetRouter())
	n.Run(router.config.GetString("service.port"))
}
