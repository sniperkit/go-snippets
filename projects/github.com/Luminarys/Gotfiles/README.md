# Gotfiles
Gotfiles is a simpler way of managing and sharing your dotfiles. Rather than create some deployment shell script, and have to constantly keep a custom directory updated with your dotfiles, gotfiles lets you do this with a few simple commands.
In addition, gotfiles makes it easy to try out dotfiles from other people by allowing you to select what repository of dotfiles you want to deploy from.

# Requirements
* Go
* rsync

# Setup
Run `make && sudo make install` and you should be good to go.

# Usage
* `gfs init [name]` creates a new dotfiles repo with a given name. This dir will be located in `~/gotfiles/[name]`
* `gfs default [name]` sets the dotfiles repo with the specified name as the default one.
* `gfs add ?[name] [file1/dir1] [file2/dir2]...` adds a file or directory to the dotfile repo with specified name. [name] is an optional argument, and if not given the default repo will be used.
* `gfs sync ?[name]` synchronizes the dotfiles in your home directory with the repo with [name]. [name] is an optional argument, and if not given the default repo will be used
* `gfs deploy -f ?[name] ?[file/dir]` deploys files from the repo to the home dir. If [name] is given, that repo is used, otherwise the default is used. 
If ?[file/dir] is given the specified file or dir will be deployed. If -f is used, files will be overwritten in the home dir, otherwise they will not be.
