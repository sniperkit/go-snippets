# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-py () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    FILENAME=$1

    MSG='\n→ Watching python (using entr)..$NC'
    BANNER='${GREEN}Output of $FILENAME:$NC'

    COMMAND="clear && echo '$BANNER' && python $FILENAME "

	FILES="find . -name \*.py"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
