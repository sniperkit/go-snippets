package metadata

import (
	"net/http"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func TagsIndexHandle(c echo.Context) error {
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

	result := s.GetTags(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}
