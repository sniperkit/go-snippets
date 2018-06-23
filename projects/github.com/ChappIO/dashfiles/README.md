# dashfiles

Manage your dotfiles with ease.

The dashfiles client allows you to install dotfiles from a git repository.

## Installation

To install the dashfiles client you only need go and git.

1. Install Git and [Go](https://golang.org/doc/install)
2. Run `go get github.com/ChappIO/dashfiles`

You are now good to go!

## Getting Started

1. If you do not have a git(hub) repository for your dotfiles yet, you should first 
   [create one](https://github.com/new).
2. To clone your repository into the dashfiles workspace you run 
   `dashfiles init <github-username>/<github-repository>` (For example: `dashfiles init ChappIO/dotfiles)`
> Note that if you do not wish to use the default setup (github using ssh) you can also pass a full git url 
> to the `init` command.
3. You can now edit you dotfiles if you wish
4. To install your current workspace to your system you simply run `dashfiles install`
