package main

/**
 * ref, https://github.com/bramvdbogaerde/go-scp
 */

import (
	"errors"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/atrox/homedir"
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

type sshConfig struct {
	Aws awsConfig `toml:"AWS"`
}

type awsConfig struct {
	PemPath        string
	UserName       string
	RemoteAddr     string
	LocalFilePath  string
	RemoteFilePath string
}

const keysConfig = "keys.toml"

var client scp.Client

func main() {

	// 1. load configuration
	cfg, err := loadSSHConfig()
	if err != nil {
		panic(err)
	}

	// 2. init client
	err = initClient(&cfg.Aws)
	if err != nil {
		panic(err)
	}

	// 3. upload file
	upload(cfg.Aws.LocalFilePath, cfg.Aws.RemoteFilePath)
	log.Printf("success upload to [%s] ==> %s", cfg.Aws.RemoteAddr, cfg.Aws.RemoteFilePath)

}

func initClient(awsCfg *awsConfig) error {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, err := auth.PrivateKey("root", awsCfg.PemPath, ssh.InsecureIgnoreHostKey())

	if err != nil {
		log.Printf("err: %v\n", err)
		return err
	}
	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod
	// Create a new SCP client
	client = scp.NewClient(awsCfg.RemoteAddr, &clientConfig)
	return nil
}

func upload(localFile, remoteFile string) error {
	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		log.Println("Couldn't establisch a connection to the remote server ", err)
		return err
	}

	// Open a file
	f, err := os.Open(localFile)
	if err != nil {
		return err
	}

	// Close session after the file has been copied
	defer client.Session.Close()
	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)

	client.CopyFile(f, remoteFile, "0655")
	return nil
}

// loadSSHConfig - load ssh config from toml file
// suppose the caller are in the same directory of this file
func loadSSHConfig() (*sshConfig, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("can not get caller's abs filepath")
	}

	cfgFilePath := path.Join(path.Dir(filename), keysConfig)

	var config sshConfig
	_, err := toml.DecodeFile(cfgFilePath, &config)

	config.Aws.PemPath, _ = homedir.Expand(config.Aws.PemPath)
	config.Aws.LocalFilePath, _ = homedir.Expand(config.Aws.LocalFilePath)

	return &config, err
}
