watch-phpunit () {
    MSG='\n\033[0;33m→ Watching php files to run the test (using python watchdog)..\033[0;0m'

    if [[ $1 ]]; then
      COMMAND="clear && vendor/bin/phpunit $1; printf \"$MSG\""
    else
      COMMAND="clear && vendor/bin/phpunit; printf \"$MSG\""
    fi

    eval $COMMAND
    eval "watchmedo shell-command \
        --patterns='*.php' \
        --recursive \
        --command='$COMMAND' \
        ."
}
