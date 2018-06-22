package config

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"git.andrewcsellers.com/acsellers/card_sharp/store"
	"github.com/BurntSushi/toml"
	"github.com/acsellers/multitemplate"
	"github.com/acsellers/multitemplate/helpers"
	_ "github.com/acsellers/multitemplate/terse"
	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

var (
	Conn     *store.Conn
	Dev      = flag.Bool("dev", false, "Development mode, load local config")
	ConfPath = flag.String("conf", "/etc/card_party/solo.conf", "Production Config file location")
	Config   Settings
	Cookie   *securecookie.SecureCookie
	Tmpl     *multitemplate.Template
)

type Settings struct {
	WebPort      int
	ResourcePath string
	SQLAddr      string
	SQLType      string
}

func init() {
	flag.Parse()

	var confPath string
	if *Dev {
		if _, err := os.Stat("solo.conf"); err != nil {
			log.Fatal("Missing solo.conf for config settings")
		}
		confPath = "solo.conf"
	} else {
		confPath = *ConfPath
	}
	if _, err := toml.DecodeFile(confPath, &Config); err != nil {
		log.Fatal("Parse Config File", err)
	}
	var err error
	if Conn, err = store.Open(Config.SQLType, Config.SQLAddr); err != nil {
		log.Fatal("Open SQL Connection", Config, err)
	}
	Conn.Setup()

	Cookie = securecookie.New(
		securecookie.GenerateRandomKey(32),
		securecookie.GenerateRandomKey(32),
	)

	helpers.LoadHelpers("all")
	CompileTemplates()
}

func WebPort() string {
	if Config.WebPort == 0 {
		Config.WebPort = 8100
	}
	return fmt.Sprintf(":%d", Config.WebPort)
}

func CompileTemplates() {
	Tmpl = multitemplate.New("base")
	Tmpl.Funcs(TemplateFuncs())

	Tmpl.Base = "templates"
	_, err := Tmpl.ParseGlob(filepath.Join("templates", "*.html.*"))
	if err != nil {
		log.Fatal(err)
	}
}

func TemplateFuncs() template.FuncMap {
	play_messages := []string{
		"Play This Game!",
		"Let's Get to Playing!",
		"Game On!",
		"Start Playing!",
		"Play Now!",
		"Do Want (to play)",
		"How Now Game Cow!",
		"We would like to play",
		"Start the Game Already!",
	}
	return template.FuncMap{
		"available_games": func() ([]store.Deck, error) {
			return Conn.Deck.AvailableDecks().RetrieveAll()
		},
		"play_this_game": func() string {
			return play_messages[rand.Intn(len(play_messages))]
		},
	}
}
