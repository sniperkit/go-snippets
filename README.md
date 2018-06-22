# mgit - the git multiplexer

Work on multiple git repositories at the same time.

### example usage

```bash
$ mgit clone git@github.com/phayes/* # check out all phayes repositories
$ mgit checkout testing # Check out testing branch on all repositories
```

### installation

```bash
$ wget "https://phayes.github.io/bin/current/mgit/linux/mgit.gz" #replace `linux` with `mac` for MacOSX version.
$ gunzip mgit.gz
$ sudo cp mgit /usr/bin     
$ sudo chmod a+x /usr/bin
$ echo "export GITHUB_API_TOKEN=MyGithubApiKeyGoesHere" >> ~/.bash_profile && source ~/.bash_profile # for github integration
```

### github integration

To clone multiple private github repositories using wildcare (`*`) notation, you need provide mgit an API token. 

1. Visit https://github.com/settings/applications
2. Under `Personal access tokens` click `Generate New Token` and give it the `repo` permission
3. Copy the api token provided
4. Edit your bash profile and add `export GITHUB_API_TOKEN=abc123` but replace `abc123` with the token
5. You will need to restart your bash session for settings to take effect. Alternatively run `source ~/.bash_profile`
 
