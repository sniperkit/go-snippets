package sqlite

import (
	"path"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/api/storage/driver"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(config string) (driver.Storager, error) {
	debug := false

	dbFilePath := path.Join(config, "data.db")

	db, err := gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	db.LogMode(debug)

	db.AutoMigrate(&ServerConfig{}, &Data{}, Profile{}, IP{}, Interface{}, Host{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	c.DB.Close()
}
