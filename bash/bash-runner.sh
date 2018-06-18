# load script from base dir and run it.
# `BASH_RUNNER_BASE_DIR` must be in path.
bashy () {
	if [ -z "$1" ]; then
		ls $BASH_RUNNER_BASE_DIR
	else
    	source "$BASH_RUNNER_BASE_DIR/$1.sh" && $@
    fi
}