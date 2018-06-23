# pet-plus : CLI Snippet Manager

Multi command-line snippet manager, written in Go.

# Abstract

`pet` is written in Go, and therefore you can just grab the binary releases and drop it in your $PATH.

`pet` is a simple command-line snippet manager (base on [pet](https://github.com/knqyf263/pet)).

Support change work snippet, multi snippet manager.

Add `use|rm` function to change work snippet.

Change Gist sync to AliyunOSS Bucket.

# Main features
`pet` has the following features.

- Register your command snippets easily.
- Use variables in snippets.
- Search snippets interactively.
- Run snippets directly.
- Edit snippets easily (config is just a TOML file).
- Multi snippets file management.
- Sync snippets via AliyunOSS.

# Examples
Some examples are shown below.

## Register the previous command easily

### bash/zsh
By adding the following config to `.bashrc` or `.zshrc`, you can easily register the previous command.

```
$ cat .zshrc
function prev() {
  PREV=$(fc -lrn | head -n 1)
  sh -c "pet new `printf %q "$PREV"`"
}
```

### fish
See below for details.

<img src="doc/pet02.gif" width="700">

## Select snippets at the current line (like C-r)

### bash
By adding the following config to `.bashrc`, you can search snippets and output on the shell.

```
$ cat .bashrc
function pet-select() {
  BUFFER=$(pet search --query "$READLINE_LINE")
  READLINE_LINE=$BUFFER
  READLINE_POINT=${#BUFFER}
}
bind -x '"\C-x\C-r": pet-select'
```

### zsh

```
$ cat .zshrc
function pet-select() {
  BUFFER=$(pet search --query "$LBUFFER")
  CURSOR=$#BUFFER
  zle redisplay
}
zle -N pet-select
stty -ixon
bindkey '^s' pet-select
```

### fish
See below for details.

<img src="doc/pet03.gif" width="700">


## Copy snippets to clipboard
By using `pbcopy` on OS X, you can copy snippets to clipboard.

<img src="doc/pet06.gif" width="700">

# Features

## Edit snippets

```
$ pet configure
```

## Change snippets

```
$ pet use [new_snippet_file_name.toml]
```

## Sync snippets

```
$ pet sync
```


# Usage

```
pet - Simple command-line snippet manager.

Usage:
  pet [command]

Available Commands:
  configure   Edit config file
  edit        Edit snippet file
  exec        Run the selected commands
  help        Help about any command
  list        Show all snippets
  new         Create a new snippet
  rm          Remove current snippet
  search      Search snippets
  sync        Sync snippets
  use         Change/Create the work snippet
  version     Print the version number

Flags:
      --config string   config file (default is $HOME/.config/pet/config.toml)
      --debug           debug mode

Use "pet [command] --help" for more information about a command.
```

# Snippet
Run `pet edit`  
You can also register the output of command (but cannot search).

```
[[snippets]]
  description = "echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout"
  command = "Show expiration date of SSL certificate"
  output = """
notBefore=Nov  3 00:00:00 2015 GMT
notAfter=Nov 28 12:00:00 2018 GMT"""
```

Run `pet list`

```
Description: echo | openssl s_client -connect example.com:443 2>/dev/null |openssl x509 -dates -noout
    Command: Show expiration date of SSL certificate
     Output: notBefore=Nov  3 00:00:00 2015 GMT
             notAfter=Nov 28 12:00:00 2018 GMT
------------------------------
```


# Configuration

Run `pet configure`

```
[General]
  snippetfile = "path/to/snippet" # specify snippet directory
  editor = "vim"                  # your favorite text editor
  column = 40                     # column size for list command
  selectcmd = "peco"              # selector command for edit command (peco or fzf)

[AliOSS]
  access_id = ""
  access_key = ""
  bucket_name = ""
  endpoint = ""
```

## Selector option
Example1: Change layout (bottom up)

```
$ pet configure
[General]
...
  selectcmd = "peco --layout=bottom-up"
...
```

Example2: Enable colorized output
```
$ pet configure
[General]
...
  selectcmd = "fzf --ansi"
...
$ pet search --color
```

## Tag
You can use tags (delimiter: space).
```
$ pet new -t
Command> ping 8.8.8.8
Description> ping
Tag> network google
```

Or edit manually.
```
$ pet edit
[[snippets]]
  description = "ping"
  command = "ping 8.8.8.8"
  tag = ["network", "google"]
  output = ""
```

They are displayed with snippets.
```
$ pet search
[ping]: ping 8.8.8.8 #network #google
```

## Sync
You must obtain access token.
Go https://oss.console.aliyun.com/index and create aliyun oss bucket, then set the configure (only need "AliOSS" scope).
Set that to `access_id` 、`access_key` 、`bucket_name` 、`endpoint` in `[AliOSS]`.

After setting, you can upload snippets to Aliyun OSS.
```
$ pet sync
Upload success
```

You can download snippets on another PC.
```
$ pet sync -f
Download success
```

# Installation
You need to install selector command ([fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco)).  
`homebrew` install `peco` automatically.

## Binary
Go to [the releases page](https://github.com/wixb50/pet/releases), find the version you want, and download the zip file. Unpack the zip file, and put the binary to somewhere you want (on UNIX-y systems, /usr/local/bin or the like). Make sure it has execution bits turned on. 

## Build

```
$ git clone https://github.com/wixb50/pet.git
$ cd pet
$ make install
```

# License
MIT

# Author
Elegance Tse
