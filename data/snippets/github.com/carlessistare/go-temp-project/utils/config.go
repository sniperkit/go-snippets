package utils

import "os"

func GetConfig() string {
	environment := os.Getenv("ENV")
	if len(environment) == 0 {
		environment = "development"
	}
	return environment
}