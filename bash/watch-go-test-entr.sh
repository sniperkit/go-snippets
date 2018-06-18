# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-go-test () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    MSG='\n→ Watching go test (using entr)..$NC'
    BANNER='${GREEN}Output of \"$*\":$NC'

    COMMAND="clear && echo '$BANNER' && go test $*"

	FILES="find . -name \*.go"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
