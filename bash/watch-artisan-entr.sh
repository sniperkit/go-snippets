# NOTE: watchdog must be installed.
# eg: `watchartisan list` will execute 
# `php artisan list` everytime there's change in `*.php` files
watch-artisan () {
	NC='\033[0;0m' # no color

    MSG='→ Watching php files to run artisan command (using entr)..$NC'

	FILES="find . -name \*.php -not -path './vendor/*'"
	COMMAND="clear && php artisan $@"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG'\""
}
