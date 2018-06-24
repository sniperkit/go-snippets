package metadata

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func NetworkInterfacesIndexHandle(c echo.Context) error {
	var result []string

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

	r := s.GetInterfaces(serverid)
	for _, i := range r {
		result = append(result, fmt.Sprintf("%s/", i))
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func NetworkInterfaceTypeIndexHandle(c echo.Context) error {
	var result []string

	itype := c.Param("type")
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

	r := s.GetInterfacesType(serverid, itype)
	for _, i := range r {
		result = append(result, fmt.Sprintf("%d/", i))
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func NetworkEnumeratedInterfaceIndexHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
	serverid := c.Get("SERVERID").(int)

	result := []string{
		"mac",
		"type",
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	i, _ := strconv.Atoi(index)

	r := s.GetEnumeratedInterface(serverid, itype, i)
	for _, i := range r {
		result = append(result, fmt.Sprintf("%s/", i))
	}

	if s.FloatingIPExists(serverid, itype, i) {
		result = append(result, "floating_ip/")
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func NetworkInterfaceMACAddressHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceMACAddress(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceTypeHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceType(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceIPv4IndexHandle(c echo.Context) error {
	result := []string{
		"address",
		"netmask",
		"gateway",
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func NetworkInterfaceIPv4AddressHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceIPv4Address(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceIPv4NetmaskHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceIPv4Netmask(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceIPv4GatewayHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceIPv4Gateway(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceFloatingIPIndexHandle(c echo.Context) error {
	result := []string{
		"address",
		"netmask",
		"gateway",
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func NetworkInterfaceFloatingIPAddressHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceFloatingIPAddress(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceFloatingIPNetmaskHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceFloatingIPNetmask(serverid, itype, i)

	return c.String(http.StatusOK, result)
}

func NetworkInterfaceFloatingIPGatewayHandle(c echo.Context) error {
	itype := c.Param("type")
	index := c.Param("num")
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

	i, _ := strconv.Atoi(index)

	result := s.GetInterfaceFloatingIPGateway(serverid, itype, i)

	return c.String(http.StatusOK, result)
}
