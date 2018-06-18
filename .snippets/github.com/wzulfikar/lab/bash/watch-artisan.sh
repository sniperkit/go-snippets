# NOTE: watchdog must be installed.
# eg: `watchartisan list` will execute 
# `php artisan list` everytime there's change in `*.php` files
watch-artisan () {
    eval "watchmedo shell-command \
        --patterns='*.php' \
        --recursive \
        --command='php artisan $1' \
        ."
}
