run-cpp () {
    FILENAME=$1
    eval "g++ -o $FILENAME $FILENAME.cpp && ./$FILENAME"
}

# `entr` requires unix `entr` utility
# → http://entrproject.org
watch-cpp () {
	NC='\033[0;0m' # no color
	GREEN='\033[0;32m'

    FILENAME=$1

    MSG='\n→ Watching cpp files (using entr)..$NC'
    BANNER='${GREEN}Output of $FILENAME:$NC'

    COMMAND="clear && g++ -o $FILENAME $FILENAME.cpp && echo '$BANNER' && ./$FILENAME"

	FILES="find . -name \*.cpp"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG';\""
}
