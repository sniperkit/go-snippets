---
layout: post
title: "Draft - Vim manual setup"
tagline: "Supporting tagline"
category : notes
tags : [draft]
---

## TL;DR

* <https://github.com/fabic/dude.vim>

* <http://marcgg.com/blog/2016/03/01/vimrc-example/>
* <https://statico.github.io/vim.html>
* <https://statico.github.io/vim2.html>

### Vim tips


#### Normal mode

* `ZZ`
* `ZQ`
* `gv` will re-select your last visual selection.


#### Command mode

* `<C-R> "` paste yanked text.
    * See [How to paste yanked text into Vim command line? (Stackoverflow)](http://stackoverflow.com/a/3997110/643087)

* You can use `<C-R><C-W>` on the ':' command line to paste the word under cursor
  ([wikia](http://vim.wikia.com/wiki/Word_under_cursor_for_command)).


#### Ex mode

* [About Vim Ex mode](https://en.wikibooks.org/wiki/Learning_the_vi_Editor/Vim/Modes#Ex-mode)


### ...

    $ mkdir -p dude.vim/dot_vim/bundle/

    $ cd dude.vim/

    $ touch dot_vimrc

Symlink `.vim/` and `.vimrc` :

    $ ( cd ~ && ln -sni ${OLDPWD#$HOME/}/vim/ .vim && ls -lad .vim )
    lrwxrwxrwx 1 fabi fabi 21 Feb 13 13:44 .vim -> dev/dude.vim/vim/

    (13:48:26) [fabi@winterfell] 0 ~/dev/dude.vim $ touch vimrc

    (13:48:50) [fabi@winterfell] 0 ~/dev/dude.vim $ ( cd ~ && ln -si ${OLDPWD#$HOME/}/vimrc .vimrc && ls -lad .vimrc )

Head to <https://github.com/VundleVim/Vundle.vim> :

    git clone https://github.com/VundleVim/Vundle.vim.git vim/bundle/Vundle.vim

### vimrc

    $ cat >> dot_vimrc <<EOF
    > set nocompatible              " be iMproved, required
    > filetype off                  " required
    >
    > " set the runtime path to include Vundle and initialize
    > set rtp+=~/.vim/bundle/Vundle.vim
    > call vundle#begin()
    > " alternatively, pass a path where Vundle should install plugins
    > "call vundle#begin('~/some/path/here')
    >
    > " let Vundle manage Vundle, required
    > Plugin 'VundleVim/Vundle.vim'
    >
    > " The following are examples of different formats supported.
    > " Keep Plugin commands between vundle#begin/end.
    > " plugin on GitHub repo
    > Plugin 'tpope/vim-fugitive'
    > " plugin from http://vim-scripts.org/vim/scripts.html
    > Plugin 'L9'
    > " Git plugin not hosted on GitHub
    > Plugin 'git://git.wincent.com/command-t.git'
    > " git repos on your local machine (i.e. when working on your own plugin)
    > Plugin 'file:///home/gmarik/path/to/plugin'
    > " The sparkup vim script is in a subdirectory of this repo called vim.
    > " Pass the path to set the runtimepath properly.
    > Plugin 'rstacruz/sparkup', {'rtp': 'vim/'}
    > " Install L9 and avoid a Naming conflict if you've already installed a
    > " different version somewhere else.
    > Plugin 'ascenator/L9', {'name': 'newL9'}
    >
    > " All of your Plugins must be added before the following line
    > call vundle#end()            " required
    > filetype plugin indent on    " required
    > " To ignore plugin indent changes, instead use:
    > "filetype plugin on
    > "
    > " Brief help
    > " :PluginList       - lists configured plugins
    > " :PluginInstall    - installs plugins; append `!` to update or just :PluginUpdate
    > " :PluginSearch foo - searches for foo; append `!` to refresh local cache
    > " :PluginClean      - confirms removal of unused plugins; append `!` to auto-approve removal
    > "
    > " see :h vundle for more details or wiki for FAQ
    > " Put your non-Plugin stuff after this line
    > EOF

### ...

    $ mkdir vim/colors

    $ for theme in Tomorrow{,-Night{,-Blue,-Bright,-Eighties}}.vim ; do
        echo "> $theme"
        wget -O vim/colors/$theme \
          "https://github.com/chriskempson/tomorrow-theme/raw/master/vim/colors/$theme" \
            || break;
      done



__EOF__
