watch-phpunit () {
	NC='\033[0;0m' # no color

    MSG='\n→ Watching php files to run the test (using entr)..$NC'

	FILES="find . -name \*.php -not -path './vendor/*'"

    if [[ $1 ]]; then
      COMMAND="clear && vendor/bin/phpunit $1"
    else
      COMMAND="clear && vendor/bin/phpunit"
    fi

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG'\""
}
