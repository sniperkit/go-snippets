# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-js () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    FILENAME=$1

    MSG='\n→ Watching js (using entr)..$NC'
    BANNER='${GREEN}Output of $FILENAME:$NC'

    COMMAND="clear && echo '$BANNER' && node $FILENAME"

	FILES="find . -name \*.js"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
