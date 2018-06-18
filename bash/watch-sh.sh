# eg. `watch-sh '*.php' 'say ha'`
watch-sh () {
	COMMAND=$2
	PATTERN=$1
   
   echo "Watching changes on '$PATTERN' to trigger '$COMMAND'.\nPress ctrl+c to exit."

   eval "watchmedo shell-command \
      --patterns='$PATTERN' \
      --recursive \
      --command='$COMMAND' \
      ."
}
