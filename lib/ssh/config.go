package ssh

import (
	"time"

	"github.com/gliderlabs/com/config"
)

type Config struct {
	// port to bind on
	ListenAddr string `mapstructure:"listen_addr"`

	// private key for host verification
	HostKeyPEM string `mapstructure:"hostkey_pem"`

	// maximum duration of session in minutes
	MaxTimeout time.Duration `mapstructure:"max_timeout"`

	// idle timeout for sessions in minutes
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

func (com *Component) InitializeConfig(config config.Settings) error {
	config.SetDefault("listen_addr", "127.0.0.1:2222")
	config.SetDefault("hostkey_pem", "lib/ssh/data/dev_host")
	return config.Load(&(com.Config))
}
