GitHub-Bot
==========
GitHub-Bot was created and maintained by [Nyah Check](https://github.com/Ch3ck), and it's a GitHub bot to follow users on GitHub. It works by taken a user token for any GitHub user, retrieves the list of it's followers and follows them. You can modify the code however to unfollow everyone you follow on GitHub. It uses [Go-GitHub](github.com/google/go-github/github) library to authenticate with the GitHub API and was inspired by [Jessica Frazelle](https://github.com/jessfraz)'s [ghb0t](https://github.com/jessfraz/ghb0t).

## Installation

* [Go version 1.8](https://github.com/golang/go/releases/tag/go1.8.3)

Clone Git repo:

```
$ git clone git@github.com:Ch3ck/github-bot.git
$ cd github-bot
$ go get golang.org/x/oauth2/...
$ go get github.com/Sirupsen/logrus/...
$ go get github.com/google/go-github/...

```

## Build & Run

```
$ make
```

## Usage

Create a `GitHub` token which you will use in your application.

```
$ github-bot -h
github-bot - v1.0
  -d    run in debug mode
  -seconds int
        seconds to wait before checking for new events (default 30)
  -token string
        GitHub API token
  -user string
  		GitHub username
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

## License

GitHub-Bot is licensed under [The GNU GPL License (GNU)](LICENSE).
