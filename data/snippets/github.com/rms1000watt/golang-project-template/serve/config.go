package serve

// Config is the generic configuration for the server
type Config struct {
	Host string
	Port int

	CertsPath string
	CertName  string
	KeyName   string
}
