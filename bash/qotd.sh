qotd () {
	# replace this path to your own
	d=~/code/github/lab/bash/qotd-words.txt

	rand=$[($RANDOM % `wc -l $d|sed "s/[^0-9]//g"`)+1]
	
	# `printf` ansi escape code for green 
	# so the qotd will be printed in green color.
	# for other escape code, see https://en.wikipedia.org/wiki/ANSI_escape_code
	printf '\033[0;32m'
	sed $rand"q;d" $d
}
eval "qotd"