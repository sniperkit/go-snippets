package plan

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Clone clones a plan repository to the plans path
func Clone(repoURL string) error {
	reg, err := regexp.Compile(`(.*/)|(\.git)`)
	if err != nil {
		return err
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return fmt.Errorf("Couldn't parse %s: %s", repoURL, err)
	}

	repoName := reg.ReplaceAllString(repoURL, "")
	// TODO use viper to get plansPath
	clonePath := fmt.Sprintf("%s/.orbi/plans/%s", os.Getenv("HOME"), repoName)

	switch u.Scheme {
	default:
		return fmt.Errorf("URL scheme \"%s\": for %s is not supported", u.Scheme, repoURL)
	/*
	 Check for empty string since a common url for cloning a git repo is
	 the shorter scp-like syntax url (i.e user@server:project.git).
	 This url is schemeless so u.Scheme will return an empty string for this
	 kind of url.
	*/
	case "":
		// Pattern to match the shorter url scp-like syntax for the SSH protocol
		sshPattern := `^(\w*@)?[A-Za-z0-9-._]+:[A-Za-z0-9-._]+/`

		reg, err := regexp.Compile(sshPattern)
		if err != nil {
			return err
		}

		isSSH := reg.MatchString(repoURL)

		if isSSH {
			if err := sshClone(repoURL, clonePath); err != nil {
				return err
			}
		}
	case "ssh":
		if err := sshClone(repoURL, clonePath); err != nil {
			return err
		}
	case "https":
		if err := httpsClone(repoURL, clonePath); err != nil {
			return err
		}
	}

	return nil
}

func sshClone(repoURL, clonePath string) error {
	/*
		TODO when a new release of go-git be available use this code to get an authMethod.

		auth, err := gitssh.NewPublicKeysFromFile("user", "pub_path")
		if err != nil {
			return err
		}
	*/

	keyPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("Couldn't read private key file: %s", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("Couldn't parse private key: %s", err)
	}

	_, err = git.PlainClone(clonePath, false, &git.CloneOptions{
		Auth: &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
		},
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("Couldn't get plan: %s", err)
	}

	return nil
}

func httpsClone(repoURL, clonePath string) error {
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("Couldn't get plan: %s", err)
	}

	return nil
}
