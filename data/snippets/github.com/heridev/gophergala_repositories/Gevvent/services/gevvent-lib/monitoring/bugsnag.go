package monitoring

import (
	"github.com/asim/go-micro/store"
	"github.com/bugsnag/bugsnag-go"
	log "github.com/cihub/seelog"
)

func Init(serviceName string) {
	initBugsnag(serviceName)
}

func initBugsnag(serviceName string) {
	var apikey, releaseStage string

	item, err := store.Get("monitoring/" + serviceName + "/bugsnag/apikey")
	if err != nil {
		log.Warnf("Error loading bugsnag API key, %s", err.Error())
	} else {
		apikey = string(item.Value())
	}
	item, err = store.Get("monitoring/" + serviceName + "/bugsnag/releasestage")
	if err != nil {
		releaseStage = "production"
	} else {
		releaseStage = string(item.Value())
	}

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       apikey,
		ReleaseStage: releaseStage,
	})
}
