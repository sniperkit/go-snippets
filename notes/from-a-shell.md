---
#layout: default
layout: page
title: "Shell(s): notes"
tagline: "Things one does from a shell"
category : notes
tags : [draft, wide, bash]
published: true
---

* Your [~/.bash-it/](https://github.com/fabic/bash-it) Bash shell configuration.
    - [~/.bash\_profile](https://github.com/fabic/bash-it/blob/master/dot_bash_profile) which does little but source your `~/.bashrc`.
    - [~/.bashrc](https://github.com/fabic/bash-it/blob/master/dot_bashrc) : now a big script where things happen.
    - [custom/](https://github.com/fabic/bash-it/tree/master/custom) : custom scripts
      you wrote and others you found out there.
    - [fabic/](https://github.com/fabic/bash-it/tree/master/fabic) : other personal stuff,
      of which your [~/bin/](https://github.com/fabic/bash-it/tree/master/fabic/bin) dir.
    - __TMUX (!) :__ [~/.tmux.conf](https://github.com/fabic/bash-it/blob/master/fabic/tmux.conf)
      (and that old [~/.screenrc](https://github.com/fabic/bash-it/blob/master/dot_screenrc) you got
       from _master bastian_ (Hervé was it?))
    - Your [Git shell aliases](https://github.com/fabic/bash-it/blob/master/aliases/available/git_fabic.aliases.bash)

## TMUX

From [SO](https://unix.stackexchange.com/a/373692), about that strange _mouse
right button click_ behaviour :

> -m and -M are used to set and clear the marked pane. There is one marked pane
> at a time, setting a new marked pane clears the last. The marked pane is the
> default target for -s to join-pane, swap-pane and swap-window.
> [See `man tmux`]

```bash
# /usr/bin/env bash
set -euo pipefail

# Make three vertically split windows with text in each.
tmux split-window -v
tmux split-window -v
tmux select-layout even-vertical
tmux send-keys -t 0 'echo pane zero' C-m
tmux send-keys -t 1 'echo pane one' C-m
tmux send-keys -t 2 'echo pane two' C-m

# You can now swap the current pane with an explicitly targeted pane. Here, we
# change pane ordering from 0-1-2 to 1-0-2, and back again.
tmux select-pane -t 0
tmux swap-pane -t 1
tmux swap-pane -t 1

# You can also swap panes by "marking" one and letting the target of the swap be
# implicit. Here, we change ordering from 0-1-2 to 1-0-2, and back again.
tmux select-pane -t 0
tmux select-pane -t 1 -m
tmux swap-pane
tmux swap-pane
```

__TODO:__ find out how you'd further automatize your shell xp. with bash scripts (?)

__EOF__
