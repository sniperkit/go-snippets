[Back to README](../README.md)

### Basic Commands
`medkit install dotfiles`: the basic, most often run command.  It will redot your symlinks, run your install.sh files, source your path files, etc.  If any bundles are configured in ~/.medkit, it will process them as well.

`medkit show config`: Will list out your current configuration settings, and potentially describe any validity issues. //TODO pick a convention for list, show, or get

`medkit show dotfiles`: will list all of the dotfiles that have been processed in your current environment.

### Dotfiles Management Commands
`medkit import dotfiles <file1,file2,file*,...>`: Will import the specified list of files into the users' dotfiles repo, either in the root, or optionally into the specified bundle

`medkit import bundle <bundle-name> <file1,file2,file*,...>`: Will import the specified list of files into the users' dotfiles repo, either in the root, or optionally into the specified bundle

`medkit init dotfiles`: will init a new dotfiles directory in the configured or default location (~/dotfiles) if one doesn't exist already.

`medkit init bundle <bundle-name>`: will init a new bundles directory under the configured or default dotfiles/bundles directory.

