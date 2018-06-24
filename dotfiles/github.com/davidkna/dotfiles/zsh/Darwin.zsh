#!/usr/bin/env zsh

export PATH="$PATH:/usr/local/sbin:/usr/local/share/pypy:${HOME}/.cargo/bin"

alias update='sudo softwareupdate -i -a; brew upgrade; brew cleanup; npm update -g; sudo gem update --system; gem update'

if (( $+commands[grc] )) && (( $+commands[brew] ))
then
  source "$(brew --prefix)/etc/grc.bashrc"
fi

# `o` with no arguments opens the current directory, otherwise opens the given
# location
function o() {
	if [ $# -eq 0 ]; then
		open .;
	else
		open "$@";
	fi;
}

# Empty the Trash on all mounted volumes and the main HDD.
# Also, clear Apple’s System Logs to improve shell startup speed.
# Finally, clear download history from quarantine. https://mths.be/bum
alias emptytrash="sudo rm -rfv /Volumes/*/.Trashes; sudo rm -rfv ~/.Trash; sudo rm -rfv /private/var/log/asl/*.asl; sqlite3 ~/Library/Preferences/com.apple.LaunchServices.QuarantineEventsV* 'delete from LSQuarantineEvent'"


test -e "${HOME}/.iterm2_shell_integration.zsh" && source "${HOME}/.iterm2_shell_integration.zsh"

if brew command command-not-found-init > /dev/null 2>&1; then eval "$(brew command-not-found-init)"; fi
