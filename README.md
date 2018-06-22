# ansible-docker
Dockerfile for ansible in a docker on Alpine Linux

[![](http://badge-imagelayers.iron.io/generik/ansible:latest.svg)](http://imagelayers.iron.io/?images=generik/ansible:latest 'Get your own badge on imagelayers.iron.io')

## Usage

Simple run : 

     docker run -it --rm=true generik/ansible ansible --help

Here I'm runing a playbook from my current working directory inside docker container: 

     docker run -it -v $SSH_AUTH_SOCK:/tmp/ssh.sck -e SSH_AUTH_SOCK=/tmp/ssh.sck -v `pwd`:/ansible --rm=true generik/ansible ansible-playbook myplaybook.yml


I like to forward my SSH agent inside, and mount the current working dir inside.
You can also pass AWS access keys ... 
