package metadata

import (
	"net/http"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func DnsIndexHandle(c echo.Context) error {
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

	result := s.GetDNSIndex(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func DnsNameserversHandle(c echo.Context) error {
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

	result := s.GetDNSNameservers(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func DnsSearchDomainsHandle(c echo.Context) error {
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

	result := s.GetDNSSearchDomains(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func DnsOptionsHandle(c echo.Context) error {
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

	result := s.GetDNSOptions(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}
