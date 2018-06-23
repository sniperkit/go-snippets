#
# Integrates zsh-autosuggestions into Prezto.
#
# Authors:
#   Sorin Ionescu <sorin.ionescu@gmail.com>
#

# Load dependencies.
pmodload 'editor'

# Source module files.
source "${0:h}/external/zsh-autosuggestions.zsh" || return 1

#
# Highlighting
#

# Set highlight color, default 'fg=8'.
zstyle -s ':prezto:module:autosuggestions:color' found \
  'ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE' || ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE='fg=242'

zstyle -s ':prezto:module:autosuggestions:accept' found \
  'ZSH_AUTOSUGGEST_ACCEPT_WIDGETS' || ZSH_AUTOSUGGEST_ACCEPT_WIDGETS=(
    forward-char
    end-of-line
    vi-end-of-line
  )

zstyle -s ':prezto:module:autosuggestions:accept_partial' found \
  'ZSH_AUTOSUGGEST_PARTIAL_ACCEPT_WIDGETS' || ZSH_AUTOSUGGEST_PARTIAL_ACCEPT_WIDGETS=()


# Disable highlighting.
if ! zstyle -t ':prezto:module:autosuggestions' color; then
  ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE=''
fi

#
# Key Bindings
#

if [[ -n "$key_info" ]]; then
  # vi
  bindkey -M viins "$key_info[Control]F" vi-forward-word
  bindkey -M viins "$key_info[Control]E" vi-add-eol
fi
