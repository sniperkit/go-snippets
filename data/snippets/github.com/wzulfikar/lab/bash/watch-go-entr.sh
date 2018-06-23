# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-go () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    MSG='\n→ Watching go (using entr)..$NC'
    BANNER='${GREEN}Output of \"$*\":$NC'

    COMMAND="clear && echo '$BANNER' && go run $1"

	if [ -z "$2" ]; then
		FILES="find . -name \*.go"
	else
		# allow passing custom directory to watch
		FILES="ls *.go $2/*.go"
    fi

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
