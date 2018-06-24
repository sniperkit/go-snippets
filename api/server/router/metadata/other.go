package metadata

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func IDHandle(c echo.Context) error {
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

	result := s.GetID(serverid)

	return c.String(http.StatusOK, strconv.Itoa(result))
}

func HostnameHandle(c echo.Context) error {
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

	result := s.GetHostname(serverid)

	return c.String(http.StatusOK, result)
}

func FQDNHandle(c echo.Context) error {
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

	result := s.GetFQDN(serverid)

	return c.String(http.StatusOK, result)
}

func UserDataHandle(c echo.Context) error {
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

	result := s.GetUserData(serverid)

	return c.String(http.StatusOK, result)
}

func VendorDataHandle(c echo.Context) error {
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

	result := s.GetVendorData(serverid)

	return c.String(http.StatusOK, result)
}

func PublicKeysHandle(c echo.Context) error {
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

	result := s.GetPublicKeys(serverid)

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func RegionHandle(c echo.Context) error {
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

	result := s.GetRegion(serverid)

	return c.String(http.StatusOK, result)
}
