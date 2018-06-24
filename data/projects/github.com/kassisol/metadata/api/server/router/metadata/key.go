package metadata

import (
	"net/http"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func KeysIndexHandle(c echo.Context) error {
	serverid := c.Get("SERVERID").(int)

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	result := s.GetKeys(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func KeyNameHandle(c echo.Context) error {
	serverid := c.Get("SERVERID").(int)
	key := c.Param("name")

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	result := s.GetKey(serverid, key)

	return c.String(http.StatusOK, result)
}
