package main

import (
  "go-temp-project/utils"
  "go-temp-project/model"
)

func main() {

  utils.KafkaInit()

  kafkaMonitoringLog := model.KafkaMonitoringLog{
    User: "3314ce44-3202-11e8-b467-0ed5f89f718b",
    Name: "Bob",
  }

  utils.KafkaSend("users", kafkaMonitoringLog)

}
