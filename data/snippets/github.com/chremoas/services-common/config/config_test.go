package config

import (
	"testing"
)

func TestConfiguration_Load(t *testing.T) {
	conf := Configuration{}
	if err := conf.Load("application.dist.yaml"); err != nil {
		t.Errorf("Error from Load(): %s", err)
	}

	if conf.Namespace != "your.namespace.to.register" {
		t.Error("Namespace unset")
	}
	if conf.Database.Driver != "mysql" {
		t.Error("Database.Driver unset")
	}
	if conf.Database.Protocol != "tcp" {
		t.Error("Database.Protocol unset")
	}
	if conf.Database.Host != "hostname" {
		t.Error("Database.Host unset")
	}
	if conf.Database.Port != 3306 {
		t.Error("Database.Port unset")
	}
	if conf.Database.Database != "database" {
		t.Error("Database.Database unset")
	}
	if conf.Database.Username != "username" {
		t.Error("Database.Username unset")
	}
	if conf.Database.Password != "password" {
		t.Error("Database.Password unset")
	}
	if conf.Database.Options != "options" {
		t.Error("Database.Options unset")
	}
	if conf.Database.MaxConnections != 5 {
		t.Error("Database.MaxConnections unset")
	}
	if conf.OAuth.ClientId != "client_id" {
		t.Error("OAuth.ClientId unset")
	}
	if conf.OAuth.ClientSecret != "client_secret" {
		t.Error("OAuth.ClientSecret unset")
	}
	if conf.OAuth.CallBackProtocol != "https" {
		t.Error("OAuth.CallBackProtocol unset")
	}
	if conf.OAuth.CallBackHost != "callback_host" {
		t.Error("OAuth.CallBackHost unset")
	}
	if conf.OAuth.CallBackUrl != "callback_url" {
		t.Error("OAuth.CallBackUrl unset")
	}
	if conf.Net.ListenHost != "localhost" {
		t.Error("Net.ListenHost unset")
	}
	if conf.Net.ListenPort != 80 {
		t.Error("Net.ListenPort unset")
	}
	if conf.Bot.DiscordServerId != "https://support.discordapp.com/hc/en-us/articles/206346498-Where-can-I-find-my-server-ID-" {
		t.Error("Bot.DiscordServerId unset")
	}
	if conf.Bot.BotToken != "You'r bot token here, do not prepend Bot... we'll do that for you" {
		t.Error("Bot.BotToken unset")
	}
	if conf.Discord.InviteUrl != "url to be used for invitations" {
		t.Error("Discord.InviteUrl unset")
	}
	if conf.Registry.Hostname != "localhost" {
		t.Error("Registry.Hostname unset")
	}
	if conf.Registry.Port == 0 {
		t.Error("Registry.Port unset")
	}
	if conf.Registry.RegisterTTL == 0 {
		t.Error("Registry.RegisterTTL unset")
	}
	if conf.Registry.RegisterInterval == 0 {
		t.Error("Registry.RegisterInterval unset")
	}
	if len(conf.Inputs) != 2 {
		t.Error("Inputs unset")
	}
	if conf.Chat.Slack.Token != "123456" {
		t.Error("Chat.Slack.Token unset")
	}
	if !conf.Chat.Slack.Debug {
		t.Error("Chat.Slack.Debug unset")
	}
	if conf.Chat.Discord.Token != "Bot 123456" {
		t.Error("Chat.Discord.Token unset")
	}
	if len(conf.Chat.Discord.WhiteList) != 3 {
		t.Error("Chat.Discord.WhiteList unset")
	}
	if conf.Chat.Discord.Prefix != "!" {
		t.Error("Chat.Discord.WhiteList unset")
	}
}

func TestConfiguration_Load_NoFile(t *testing.T) {
	conf := Configuration{}
	if err := conf.Load("application.nofile.yaml"); err == nil {
		t.Error("No error from Load() with no file")
	}

	if conf.initialized {
		t.Error("conf failed to load but said it was initialized.")
	}
}

func TestConfiguration_Load_InvalidFile(t *testing.T) {
	conf := Configuration{}
	if err := conf.Load("application.invalid.yaml"); err == nil {
		t.Error("No error from Load() with invalid file")
	}

	if conf.initialized {
		t.Error("conf failed to load but said it was initialized.")
	}
}

func TestConfiguration_Load_NoNamespace(t *testing.T) {
	conf := Configuration{}
	if err := conf.Load("application.noname.yaml"); err != nil {
		t.Errorf("Error from Load(): %s", err)
	}

	if conf.Namespace == "" {
		t.Error("Namespace not set to a default")
	}
}

func TestConfiguration_IsInitialized_NotInitialized(t *testing.T) {
	conf := Configuration{}

	if conf.IsInitialized() {
		t.Error("Configuration was initialized after construction?  Really?")
	}
}

func TestConfiguration_IsInitialized(t *testing.T) {
	conf := Configuration{}
	conf.Load("application.dist.yaml")

	if !conf.IsInitialized() {
		t.Error("Nothing was even validated... should have been initialized!")
	}
}

func TestConfiguration_ExtensionsLoaded(t *testing.T) {
	conf := Configuration{}
	conf.Load("application.with-extensions.yaml")

	if !conf.IsInitialized() {
		t.Error("Nothing was even validated... should have been initialized!")
	}

	if conf.Extensions == nil {
		t.Error("No extensions loaded.")
	}

	if conf.Extensions["mongoDb"].(map[interface{}]interface{})["host"].(string) != "localhost" {
		t.Error("mongoDb.host was not set")
	}

	if conf.Extensions["mongoDb"].(map[interface{}]interface{})["port"] != 1234 {
		t.Error("mongoDb.port was not set")
	}
}
