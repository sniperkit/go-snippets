package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	client, err := authenticate("10.71.1.3", "tsl", "tsl")
	//client, err := authenticate("10.71.20.102", "liwei", "Lee123!@#")
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	for _, cmd := range os.Args[1:] {
		stdin.Write([]byte(cmd + "\n"))
	}

	stdin.Write([]byte("logout\n"))
	stdin.Write([]byte("N\n"))

	session.Wait()
}

func authenticate(host, username, password string) (ssh.Client, error) {
	hostKey, err := checkHostKey(host)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// Create client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		Timeout:         time.Second * 5,
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", host+":22", config)

	return *client, err
}

func checkHostKey(host string) (ssh.PublicKey, error) {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
			}
			break
		}
	}

	if hostKey == nil {
		return nil, errors.New(fmt.Sprintf("no hostkey for %s", host))
	}

	return hostKey, nil
}
