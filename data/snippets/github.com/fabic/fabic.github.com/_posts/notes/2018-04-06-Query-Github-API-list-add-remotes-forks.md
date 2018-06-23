---
layout: note
title: "Github/API: List forks of a repository + jq.sh"
tagline: "How to query the Github API so as to list all forks and add remotes for these."
published: true
category : notes
tags : [note, git, cli, API]
---

__TL;DR:__ This will query the [Github API (v3)](https://developer.github.com/v3/repos/forks/)
endpoint using Curl, and select the SSH/Git URL as well as the fork owner username
using [jq](https://stedolan.github.io/jq/), that one may process through Sed and
a Bash while-read loop for adding Git remotes :

```bash
$ curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/ulli-kroll/mt7610u/forks |
      jq '.[] | [.ssh_url, .owner.login] | @sh' |
        sed -e 's/"\(.*\)"/\1/' -e "s/'//g" |
          while read url login ;
          do
            echo "~~> Adding remote '$login' for $url :";
            git remote add "$login" "$url" || (echo "ERROR!") ;
          done
```

We resort to Sed for removing the wrapping double-quotes; as well as single-quotes.

The most simple query using Curl and basic auth/n :

```bash
$ curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/ulli-kroll/mt7610u/forks
```

See [developer.github.com/v3/](https://developer.github.com/v3/#authentication)
about the available authentication mecanisms.

Piping the output through __jq__, selecting data :

```bash
$ curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/ulli-kroll/mt7610u/forks |
      jq '.[] | {name: .full_name, url: .html_url, git: .ssh_url, homepage: .homepage, has_issues: .has_issues, has_wiki: .has_wiki}'
```

Getting just the list of Git/SSH URLs :

```bash
$ curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/ulli-kroll/mt7610u/forks |
      jq '.[] | .ssh_url'
```

Along with owner login, wrapped into an array :

```bash
curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/ulli-kroll/mt7610u/forks |
    jq '.[] | [.ssh_url, .owner.login]'
```

Filtered with `|@sh` provides an output that is very close to a table :

```bash
curl -s -u fabic -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/ulli-kroll/mt7610u/forks |
    jq '.[] | [.ssh_url, .owner.login]'
```

__NOTE:__ There's also Bash script [ok.sh](https://github.com/whiteinge/ok.sh),
but it seems not to provide a command for listing forks.

__EOF__

