#
# Bash Profile - Helpers
#

BASH_PROFILE_FILES=\
(
	".bash_profile.alias"
	".bash_profile.alias.git"
	".bash_profile.alias.go"
	".bash_profile.funcs"
	".bash_profile.funcs.git"
	".bash_profile.funcs.go"
	".bash_profile.output"
	".bash_profile.stats"
)

#### Common aliases
alias bigs="du -hm . | sort -hnr | head -n 30"
alias biggy="du -m . | sort -nr | head -n 30"
alias lsd="du -sh ."

# Function bash_parent_dir() allows to return the parent dir of bash source script called 
function bash_parent_dir() {
	local dir=${1:-"`pwd`"}
	echo $(dirname "${BASH_SOURCE[0]}")
}

function ensure_dir() {
	local prefix_path=${1:-".meta"}
	local force=${2:-"true"}
	local quiet=${3:-"true"}
	local opts=""
	[ "$force" == "true" ] && opts="-p"
	local res=$(mkdir ${opts} ${prefix_path})
	[ "$res" == "true" ] && opts="-p"
}

function ls-dirs-size {
	local limit=${1:-"50"}
	local export_file=${1:-".meta/ls-dirs.by_size.top${limit}.output"}
	ensure_dir
	du -m . | sort -nr | head -n $limit > ${export_file}
}

function ff { 
	local pattern=${1:-"test"}
	find . -type f -name "$pattern" -print 
}

function ls_dirs_size {
	local start=$(timer)
	local limit=${1:-"50"}
	local export_file=${1:-".meta/ls-dirs.by_size.top${limit}.output"}
	ensure_dir
	du -m . | sort -nr | head -n $limit > ${export_file}
	echo $(timer $start $FUNCNAME)
	separator
}
