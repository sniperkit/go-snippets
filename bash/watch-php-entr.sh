# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-php () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    FILENAME=$1

    MSG='\n→ Watching php (using entr)..$NC'
    BANNER='${GREEN}Output of $FILENAME:$NC'

    COMMAND="clear && echo '$BANNER' && php $FILENAME"

	FILES="find . -name \*.php"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
