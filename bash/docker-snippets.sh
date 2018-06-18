# run ubuntu bash shell
alias ubuntu="docker run -it --rm phusion/baseimage bash"
alias centos="docker run -it --rm centos bash"

# `docker ps` format string for vertical display
export DOCKER_PS_FORMAT="ID\t{{.ID}}\nNAME\t{{.Names}}\nIMAGE\t{{.Image}}\nPORTS\t{{.Ports}}\nCOMMAND\t{{.Command}}\nCREATED\t{{.CreatedAt}}\nSTATUS\t{{.Status}}\n"

# run bash from given container name and remove the container on exit
bash-from () {
	eval "docker run -it --rm $1 bash"
}
