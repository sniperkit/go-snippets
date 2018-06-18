# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-java () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    FILENAME=$1

    MSG='\n→ Watching java (using entr)..$NC'
    BANNER='${GREEN}Output of $FILENAME.java:$NC'

    COMMAND="clear && echo '$BANNER' && javac $FILENAME.java && java $FILENAME"

	FILES="find . -name \*.java"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
