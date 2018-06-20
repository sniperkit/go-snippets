package config

import (
	"errors"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"os"
	"testing"
)

func TestConfiguration_NewService(t *testing.T) {
	os.Args = []string{os.Args[0], "--configuration_file=application.dist.yaml"}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "defaultName", errInit)
	if _, ok := service.(micro.Service); !ok {
		t.Error("service is not a micro.Service!")
	}

	err := service.Run()
	if err == nil {
		t.Error("Uh?  What?  I got not error?")
	}
}

func TestConfiguration_NewService_NoName(t *testing.T) {
	os.Args = []string{os.Args[0], "--configuration_file=application.noname.yaml"}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "defaultname", errInit)

	err := service.Run()

	if err != nil && err.Error() != errInit(nil).Error() {
		t.Errorf("Error not set to nil, message is: %s", err)
	}
}

func TestConfiguration_NewService_NoName_NoDefault(t *testing.T) {
	os.Args = []string{os.Args[0], "--configuration_file=application.noname.yaml"}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "", errInit)

	err := service.Run()

	if err == nil {
		t.Errorf("Error not set to nil, message is: %s", err)
	}
}

func TestConfiguration_NewService_NoConfLoaded(t *testing.T) {
	os.Args = []string{os.Args[0], "--configuration_file=application.derp.yaml"}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "defaultName", NilInit)

	err := service.Run()
	if err == nil {
		if _, ok := service.(micro.Service); ok {
			t.Error("Error is nil and yet we have a valid service?")
		}
		t.Error("Error set to nil, at least we don't have a valid service")
	}
}

func TestConfiguration_NewService_NilInit(t *testing.T) {
	os.Args = []string{os.Args[0], "--configuration_file=application.derp.yaml"}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "defaultName", NilInit)

	err := service.Run()
	if err == nil {
		if _, ok := service.(micro.Service); ok {
			t.Error("Error is nil and yet we have a valid service?")
		}
		t.Error("Error set to nil, at least we don't have a valid service")
	}
}

func TestNewService_WithBlankConfigurationEnvVar(t *testing.T) {
	os.Setenv("CONFIGURATION_FILE", "")
	os.Args = []string{os.Args[0]}
	cmd.DefaultCmd = cmd.NewCmd()

	service := NewService("1.1.1", "defaultType", "defaultName", NilInit)

	err := service.Run()
	if err == nil {
		t.Error("Expected an error but received nil.")
	}
}

func TestNilInit(t *testing.T) {
	err := NilInit(&Configuration{})
	if err != nil {
		t.Error("Received an error from a 1 line function that only returns nil?  What did you change and why did you change it?!?!?!  Seriously though... the name says it all... NIL INIT!! :P")
	}
}

// For testing purposes only, I want to boot the service to ensure parsing happens but have no way to programatically stop it
// at least not that I've found.
func errInit(configuration *Configuration) error {
	return errors.New("I'm supposed to fail")
}
