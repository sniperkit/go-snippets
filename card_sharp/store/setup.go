package store

import (
	"log"
	"os"

	"github.com/acsellers/dr/migrate"
)

func (c *Conn) Setup() error {
	db := migrate.Database{
		DB:         c.DB,
		Schema:     Schema,
		Translator: NewAppConfig("postgres"),
		DBMS:       migrate.Postgres,
		Log:        log.New(os.Stdout, "Migrate: ", 0),
	}
	err := db.Migrate()
	if err != nil {
		return err
	}

	if c.Account.Email().Eq("andrew@andrewcsellers.com").Count() == 0 {
		acc := Account{
			Email: "andrew@andrewcsellers.com",
		}
		acc.SetPassword(os.Getenv("ADMIN_PASS"))
		return acc.Save(c)
	}
	c.SetupBuiltinGames()
	return nil
}

func (c *Conn) SetupBuiltinGames() {
	if c.Deck.Name().Eq("Fill in the Blanks").Count() == 0 {
		SetupFillIn(c)
	}
	if c.Deck.Name().Eq("Description Roulette").Count() == 0 {
		SetupDescribe(c)
	}
}
