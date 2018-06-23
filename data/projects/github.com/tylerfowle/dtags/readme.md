# :bookmark: dtags
#### directory tagging and navigation utility.

Dtags makes tagging, and change directories fast and easy. Add arbitrary tags
to any directory, making it easy to jump back to those directories easy.

## Installation

Installation consists of 3 parts.
1. an alias in your bashrc/zshrc.
2. bash helper script that is used to launch commands and `cd` to tagged directories.
3. the guts, a go utility that does all the heavy lifting.


#### Download dtags and set permissions:
```
sudo curl -L https://github.com/tylerfowle/dtags/releases/download/v0.1.1/dtags -o /usr/local/bin/dtags
sudo chmod +x /usr/local/bin/dtags
```

#### Install bash helper script:

```
dtags install
```

#### Alias: add to your bashrc/zshrc:

```
alias dt=". ~/.config/dtags/dt"
```

## Usage
Command | Description
---     | ---
`<string>`         | returns a path from a tag, `cd`s you to the directory when called from bash helper script/alias
`add`              | add tag to current _working_ directory
`del`              | delete tag and associated path from database, dtags doesnt care were this command is ran from
`ls`               | prints all tags and their associated directories
`list`             | lists all tags on current working directory
`completion`       | returns a list of all tags in database.  (used for bash/zsh completion)
`install`          | install dt bash helper script

## Todos:
- [ ] add images/gifs to readme
- [x] make tags case insensitive
- [x] zsh completion
- [ ] bash completion
- [x] confirm overwrite?
- [x] add install script
- [ ] add uninstall script
- [x] add ability to manually enter path when adding tag
- [ ] ~~bash helper script - make dtags path easy to change~~
- [ ] add ability to add and delete multiple tags at once?
- [ ] ~~config file? alias name? install location? bucket name?~~
