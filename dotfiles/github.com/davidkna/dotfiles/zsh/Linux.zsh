#!/usr/bin/env zsh

export PATH="$PATH:$(ruby -e 'print Gem.user_dir')/bin:$HOME/.yarn/bin"

if [ -d ~/.linuxbrew ]; then
	export PATH="$HOME/.linuxbrew/bin:$PATH"
	export MANPATH="$(brew --prefix)/share/man:$MANPATH"
	export INFOPATH="$(brew --prefix)/share/info:$INFOPATH"
fi;

# `o` with no arguments opens the current directory, otherwise opens the given
# location
function o() {
	if [ $# -eq 0 ]; then
		gnome-open .;
	else
		gnome-open "$@";
	fi;
}

alias cp="cp --reflink=auto --sparse=always"
alias update='sudo dnf upgrade -y; sudo npm update -g; sudo gem update --system; gem update'
