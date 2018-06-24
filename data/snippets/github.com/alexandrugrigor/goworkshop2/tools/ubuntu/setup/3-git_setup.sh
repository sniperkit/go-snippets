#!/bin/bash

echo "Checking if git must be installed"
export git_installed=$(which git)

if [ -z "$git_installed" ]; then
    echo "Git is not installed"
    echo "Installing git. May require your password"
    sudo apt-get install git
else
    echo "Git is already installed"
fi

export workshop_dir=${HOME}/workspace/go/src/goworkshop

echo "Cloning the project from github."
echo "The project will be downloaded under ${workshop_dir}"

mkdir -p ${workshop_dir}

#clone with https since ssh requires the key to be set up
git clone https://github.com/alexandrugrigor/goworkshop.git ${workshop_dir}

#now change the remote url of the project to work with ssh instead of https
git --git-dir=${workshop_dir}/.git --work-tree=${workshop_dir} remote set-url origin git@github.com:alexandrugrigor/goworkshop.git

