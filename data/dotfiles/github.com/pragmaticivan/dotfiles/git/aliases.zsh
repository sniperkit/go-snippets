#!/bin/sh
if which hub >/dev/null 2>&1; then
	alias git='hub'
fi

alias g='git'

gi() {
	curl -s "https://www.gitignore.io/api/$*"
}
