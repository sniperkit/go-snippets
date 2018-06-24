//Filename: Main.go
//Author: Nyah Check
//Purpose: GitHub Bot to increase following and print following list.
//Token: 91b804cf541f1e923004b11e95af94249192b54c
//Licence: GNU PL 2017


package main

import (
    "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

const (
	// LOGGER is what is printed for help/info output
	LOGGER = "github-bot - %s\n"
	// VERSION is the binary version.
	VERSION = "v1.0"
)

var (
	token    string
	interval string
    kmd      string
    usr		 string
	lastChecked time.Time
    
	debug   bool
	version bool
)

type UserData struct {
    Login             string
	ID                int
	HTMLURL           string
	Location          string
	Email             string
}

func init() {
	// parse flags
	flag.StringVar(&token, "token", "", "GitHub API token")
	flag.StringVar(&usr, "user", "", "GitHub user(Must have many followers)")
	flag.StringVar(&interval, "interval", "30s", "check interval (ex. 5ms, 10s, 1m, 3h)")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(LOGGER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if token == "" {
		token = "ee66081a0288f9d3b010bada4f67ee0df277ca04"
	}
	
	if usr == "" {
		usr = "torvalds"
	}
}


func main() {
	var ticker *time.Ticker
	// On ^C, or SIGTERM handle exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			ticker.Stop()
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Create the github client.
	client := github.NewClient(tc)

	// Get the authenticated user, the empty string being passed let's the GitHub
	// API know we want ourself.
	user, _, err := client.Users.Get(usr)
	if err != nil {
		logrus.Fatal(err)
	}
	username := *user.Login
	
	// parse the duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		logrus.Fatalf("parsing %s as duration failed: %v", interval, err)
	}
	ticker = time.NewTicker(dur)

	logrus.Infof("Bot started for user %s.", username)
	//logrus.Infof("Enter GitHub user: ");
	//fmt.Scanf("%s", &usr)
	
	for range ticker.C {
	    numUsers := 50
	    pageNum := 2
	    
		if err := getFollowing(client, username, numUsers, pageNum); err != nil { //This parts work well
	       logrus.Fatal(err)
	   }
	    
	   if err := getFollowers(client, username, numUsers, pageNum); err != nil {
	       logrus.Fatal(err)
	   }
	   
	   if err := followUsers(client, usr,numUsers, pageNum); err != nil {
	        logrus.Fatal(err)
	   }
	    
	    /**
	     * Add this program to the cron jobs so it's executed every hour.
	     
	    if err := unFollow(client, username, numUsers, pageNum); err != nil {
			logrus.Fatal(err)
		} */
		
	}
}

// getFollowers iterates over all followers received for user.
func getFollowers(client *github.Client, username string, numUsers, pageNum int) error {
    opt := &github.ListOptions{
			    Page:    pageNum,
			    PerPage: numUsers,
	        }
    
    followers, resp, err := client.Users.ListFollowers(username,opt)
	if err != nil {
		return err
	}

    //writes user details to file.
    saveData("logs/followers.json", followers, pageNum)
	
	// Return early if we are on the last page.
	if pageNum == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	pageNum = resp.NextPage
	return getFollowers(client, username, numUsers, pageNum)
}


//saveData to file.
func saveData(file string, data []*github.User, pageNum int) (error) {
    var in *os.File
    var err error
    
    if pageNum == 1 {
        in, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
    } else {
        in, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
    }
    
    if err != nil {
        return err
    }
    defer in.Close()
    
    //serialize the data
    newdata :=  []UserData{}
    out, er := json.Marshal(data)
    err = json.Unmarshal(out, &newdata)
    out, er = json.Marshal(newdata)
    er = ioutil.WriteFile(file, out, 0644)
    //fmt.Fprintf(in, string(out))

    return er
}


// getFollowing iterates over the list of following and writes to file using a gob object
func getFollowing(client *github.Client, username string, numUsers, pageNum int) error {
	opt := &github.ListOptions{
			    Page:    pageNum,
			    PerPage: numUsers,
	        }
	
    following, resp, err := client.Users.ListFollowing(username, opt)//to test properly whether to parse resp instead inloop
	if err != nil {
		return err
	}

    saveData("logs/following.json", following, pageNum)//writes details to file.
	// Return early if we are on the last page.
	if pageNum == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	pageNum = resp.NextPage
	return getFollowing(client, username, numUsers, pageNum)
}


// followUsers, gets the list of followers for a particular user and followers them on GitHub.
// This requires authentication with the API.
func followUsers(client *github.Client, username string, numUsers, pageNum int) error {
    opt := &github.ListOptions{
			    Page:    pageNum,
			    PerPage: numUsers,
	        }
    //client.Users.Follow(username) //first of all follow this user.
    usrs, resp, err := client.Users.ListFollowing(username, opt) //to test properly whether to parse resp instead inloop
	if err != nil {
		return err
	}
	
	//fmt.Printf("\nAre we here yet.....\n\n")
	logrus.Infof("%s has %+v, Curr: %+v",username, resp.LastPage, pageNum)

	for _, usr := range usrs {
		//Follow user
		res, _ := client.Users.Follow(*usr.Login)
        fmt.Printf("%+v", res)
	}
	
    // Return early if we are on the last page.
	if pageNum == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	pageNum = resp.NextPage
	return followUsers(client, username, numUsers, pageNum)
}


// Unfollow all GitHub users on one's follower list.
func unFollow(client *github.Client, username string, numUsers, pageNum int) error {
	opt := &github.ListOptions{
			    Page:    pageNum,
			    PerPage: numUsers,
	        }
	
    usrs, resp, err := client.Users.ListFollowing(username, opt)
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		//Follow user
		res, e := client.Users.Unfollow(*usr.Login)
        if err != nil {
            panic(e.Error())
        }
        
        logrus.Infof("%+v", res)
	}

	// Return early if we are on the last page.
	if pageNum == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	pageNum = resp.NextPage
	return unFollow(client, username, numUsers, pageNum)
}


func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
