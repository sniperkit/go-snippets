# get length of whatever passsed in
len () {
	NC='\033[0m' # no color
	ORANGE='\033[0;33m'
	
	for arg in "$@" # `$@` represents all arguments
	do
		echo "→ Length of '$ORANGE$arg$NC' is $ORANGE${#arg}$NC"
	done
}
