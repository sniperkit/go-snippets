package main

import (
	"context"
	"net"

	"github.com/Xe/ln"
	"github.com/Xe/pawd/cmd/pawd/database"
	pawd "github.com/Xe/pawd/proto"
	"github.com/asdine/storm"
	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

type Config struct {
	DBPath string `env:"DB_PATH" envDefault:"./var/pawd.db"`
	//HyperAccessKey  string   `env:"HYPER_ACCESS_KEY,required"`
	//HyperSecretKey  string   `env:"HYPER_SECRET_KEY,required"`
	HTTPSPort       string   `env:"TLS_PORT" envDefault:"443"`
	GRPCTLSBindhost string   `env:"GRPC_TLS_BINDHOST" envDefault:"0.0.0.0:4228"`
	AdminEmails     []string `env:"ADMIN_EMAILS" envSeparator:","`
}

type Server struct {
	Config
	db *storm.DB

	us database.Users
	tk database.Tokens
}

func main() {
	ctx := context.Background()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("parsing config"))
	}

	db, err := storm.Open(cfg.DBPath)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("opening database"))
	}

	us := database.NewUsersStorm(db, cfg.AdminEmails)
	tk := database.NewTokensStorm(db, us)

	s := &Server{
		Config: cfg,
		db:     db,

		us: us,
		tk: tk,
	}

	a := &Auth{Server: s}
	gs := grpc.NewServer()

	l, err := net.Listen("tcp", cfg.GRPCTLSBindhost)
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("opening grpc bind"))
	}

	pawd.RegisterAuthServer(gs, a)

	ln.Log(ctx, ln.Action("now listening for grpc"), ln.F{"bindhost": cfg.GRPCTLSBindhost})
	gs.Serve(l)
}
