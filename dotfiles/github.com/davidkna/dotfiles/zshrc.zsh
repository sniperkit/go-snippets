#!/usr/bin/env zsh

source ~/.zsh/common.zsh

if [ -f ~/.local/share/zsh-goggles/plugins.zsh ];
then
    source ~/.local/share/zsh-goggles/plugins.zsh
fi

# Command completion
autoload -U compinit
compinit

# Case insensitive
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'
