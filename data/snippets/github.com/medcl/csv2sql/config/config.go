package config

type AppConfig struct {
	UIEnabled bool `config:"ui_enabled"`
}

var appConfig AppConfig

func GetAppConfig() AppConfig {
	return appConfig
}

func SetAppConfig(cfg AppConfig) {
	appConfig = cfg
}
