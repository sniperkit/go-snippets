package api

import (
	"github.com/dukex/mixpanel"
)

type metricsResource struct {
	MX *mixpanel.Mixpanel
}

// NewMetricsResource is an exported method to create a metric client.
func NewMetricsResource(userUUID string) *metricsResource {
	m := &metricsResource{}
	client := mixpanel.NewMixpanel("mytoken")
	client.Identify(userUUID)
	m.MX = client
	return m
}
