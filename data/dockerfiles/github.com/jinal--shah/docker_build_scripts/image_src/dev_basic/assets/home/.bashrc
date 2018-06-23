for script in /etc/profile.d/*.sh ; do
    if [ -r $script ] ; then
        . $script
    fi
done

if [[ -r /etc/profile.d/bash_prompt.sh ]]; then
    . /etc/profile.d/bash_prompt.sh
elif [[ -r /etc/profile.d/color_prompt ]]; then
    . /etc/profile.d/color_prompt
else
    export PS1='\[\033[01;32m\]\h\[\033[01;36m\]\W$ \[\033[00m\]'
fi

if [[ -d ~/profile.d ]]; then
    for f in ~/profile.d/*; do
        . $f || echo "ERROR: ... could not source $f"
    done
fi

forwardSsh_help() {

    cat <<EOM
usage:  forwardSsh
... adds all of your .ssh keys to an ssh-agent for the current shell
EOM

}

forwardSsh() {
    echo "... generating agent for ssh forwarding in cluster"
    pkill ssh-agent
    eval $(ssh-agent)
    for privateKey in $(ls -1 $HOME/.ssh/id_* | grep -v '\.pub')
    do
        ssh-add "$privateKey"
    done
    ssh-add -l # verify your key has been added to the key-ring
}

forwardSsh

alias gtree='tree -a -C -I .git'
