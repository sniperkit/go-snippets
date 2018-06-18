# Alias to directory for code stuffs (personal preference). 
# Type `code` in your terminal and it will bring you to 
# the directory where you place your codes. This way, 
# you can jump to your code directly from anywhere in 
# your terminal which can help reduce mental burden. 
code () {
	CODE_DIRECTORY="~/code"
	if [[ $1 ]]; then
    	eval "cd $CODE_DIRECTORY/$1"
	else
    	eval "cd $CODE_DIRECTORY"
  	fi
}

# edit zsh config using sublime
alias zshconfig="subl ~/.zshrc"

# update zsh config
alias zshso="source ~/.zshrc" # 

# edit hosts config using sublime text
alias hostsconfig="subl /etc/hosts"

# laravel artisan
alias artisan="php artisan" 

# python alias
alias py2=python2
alias py3=python3
alias py=python3
alias pyman="python3 manage.py" # for use with django

# alias for asciinema record command
alias rec="asciinema rec"

alias hist=history # `history` is aliased to fc -l 1

# convert video file to gif using gifgen
# and display osx notification once finished.
gif () {
	eval "gifgen $1  && terminal-notifier -title Gifgen -message \"Finished converting video to gif.\""
}

# scrape image from given url using image-scraper
# and display osx notification once finished.
scrape-image () {
	eval "image-scraper $1 && terminal-notifier -title 'Image Scraper' -message \"Finished scraping images from $1\""
}

# google search, from command line.
# see: https://github.com/jarun/googler
google () {
	eval "googler --colors Gcdgxy $1"
}

# find duplicate words in file
dupes () {
	rev $1 | cut -f1 -d/ | rev | sort | uniq -d
}

# watch ts files and execute index.ts or `$1`
# dependencies: nodemon, typescript, ts-node
nodemon-ts () {
	if [[ $1 ]]; then
		nodemon --watch '*.ts' --exec 'ts-node' $1
	else
		nodemon --watch '*.ts' --exec 'ts-node' index.ts
  	fi
}

# make symlink of current directory in $GOPATH/src
# e.g. create `$GOPATH/src/github.com/my-package`:
# `cd ~/path-to/my-package && golink github.com/wzulfikar`
golink() {
	eval "ln -s $(pwd) $GOPATH/src/$1"
}