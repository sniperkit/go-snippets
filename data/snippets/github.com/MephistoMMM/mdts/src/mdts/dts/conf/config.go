package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type confunit struct {
	// OutHostport listen requests from third via WAN
	OutHostport string
	// InHostport listen request from service via LAN
	InHostport string

	// Usehttps determinate to use https in OutHostport or not.
	Usehttps bool

	ServerCrtRoot string
	ServerCrt     string
	ServerKey     string

	ReqPoolSize int
	// second
	ReqClientTimeOut int

	EndPoints []string
	EtcdPath  string
}

type envconf struct {
	Production  confunit
	Development confunit
	Testing     confunit
}

var (
	configPath = "config.json"
)

// env
const (
	// DEV 开发环境
	DEV = iota
	// PRO 生产环境
	PRO
	// TEST 测试环境
	TEST

	envPattern       = "GODTSMODE"
	envConfigPattern = "GODTSCONF"

	devPattern  = "develop"
	proPattern  = "production"
	testPattern = "testing"
)

var (
	// Env 程序环境 默认为开发环境
	Env = DEV

	// OutHostport listen requests from third via WAN
	OutHostport string
	// InHostport listen request from service via LAN
	InHostport string

	// Usehttps 对第三方服务器使用 https
	Usehttps bool

	// ServerCrt 服务端证书地址
	ServerCrt string
	// ServerKey 服务端秘钥地址
	ServerKey string

	// ReqPoolSize 请求Client池大小
	ReqPoolSize int
	// ReqClientTimeOut 请求Client超时时间
	ReqClientTimeOut time.Duration

	// EndPoints etcd节点
	EndPoints []string
	// EtcdPath 服务的前缀
	EtcdPath string
)

var defaultconfunit = confunit{
	OutHostport: ":8081",
	InHostport:  ":8080",
	Usehttps:    false,

	ServerCrtRoot: "./tls/",
	ServerCrt:     "/server.crt",
	ServerKey:     "/server.key",

	ReqPoolSize:      10,
	ReqClientTimeOut: 2, // unit: second

	EndPoints: []string{"150.95.157.181:2385", "150.95.157.181:2381", "150.95.157.181:2383"},
	EtcdPath:  "broker/",
}

var defaultconf = envconf{defaultconfunit, defaultconfunit, defaultconfunit}

func loadConfFile(confPath string) *envconf {
	byt, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Println("Not Find Config File.", err)
		return &defaultconf
	}

	var localconf envconf

	err = json.Unmarshal(byt, &localconf)
	if err != nil {
		log.Fatalln("Config File Is Invalid.", err)
	}

	return &localconf
}

// ShowEnv show the base infos of enviroment.
func ShowEnv() {
	if Env == TEST {
		log.Println("Enviroment: TEST")
	} else if Env == DEV {

		log.Println("Enviroment: DEV")
	} else {
		log.Println("Enviroment: PRO")
	}
	log.Printf("Use Https[%t]\n", Usehttps)
}

func init() {
	if p := os.Getenv(envConfigPattern); p != "" {
		configPath = p
	}
	envconfContent := loadConfFile(configPath)

	// 根据环境变量选择当前环境的配置
	var cc confunit
	switch os.Getenv(envPattern) {
	case proPattern:
		Env = PRO
		cc = envconfContent.Production
		// 设置 gin 为发行环境
		gin.SetMode(gin.ReleaseMode)
	case testPattern:
		Env = TEST
		cc = envconfContent.Testing
	default:
		Env = DEV
		cc = envconfContent.Development
	}

	// Server
	OutHostport = cc.OutHostport
	InHostport = cc.InHostport

	Usehttps = cc.Usehttps

	// Request
	ReqClientTimeOut = time.Second * time.Duration(cc.ReqClientTimeOut)
	ReqPoolSize = cc.ReqPoolSize

	// Certificates
	ServerCrt = path.Join(cc.ServerCrtRoot, cc.ServerCrt)
	ServerKey = path.Join(cc.ServerCrtRoot, cc.ServerKey)

	// Etcd
	EndPoints = cc.EndPoints
	EtcdPath = cc.EtcdPath

}
