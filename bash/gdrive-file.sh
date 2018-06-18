# search file in google drive
# see: https://github.com/prasmussen/gdrive
gdrive-file () {
	query=$@

	# reject if query is empty
	if [ -z "$query" ]; then
	    echo "Command invalid - please provide search query for file.\nExample: gdrive-file monthly report"
	    return
	fi

	command="gdrive list --query \"name contains '$query'\" --max 20"
	echo "Executing '$command'.."
	eval $command
}