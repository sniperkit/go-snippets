# NOTE: watchdog must be installed.
# eg: `watchartisan list` will execute 
# `php artisan list` everytime there's change in `*.php` files
watch-asm () {
	NC='\033[0;0m' # no color

    MSG='→ Using entr to watch asm files..$NC'

	FILES="find . -name \*.asm"

	FILENAME=$1
	OUTPUT_FILE=asm_tmp.o
	COMMAND="clear && nasm -f macho64 -o $OUTPUT_FILE $1 && ld $OUTPUT_FILE -o $OUTPUT_FILE && ./$OUTPUT_FILE"

	eval "$FILES | entr sh -c \"$COMMAND; printf '$MSG'\""
}
