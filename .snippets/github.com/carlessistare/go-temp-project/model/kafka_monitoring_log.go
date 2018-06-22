package model

import "go-temp-project/utils"

type KafkaMonitoringLog struct {
	Environment string `json:"environment"`
	User string `json:"user"`
	Name string `json:"name"`
}

func (h KafkaMonitoringLog) Init() {
	h.Environment = utils.GetConfig()
}
