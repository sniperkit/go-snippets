#!/bin/sh
alias d='docker'

function docker_prune() {
	docker system prune --volumes -fa
}

function drmc() {
  docker service rm $(docker service ls -q)
  docker rm -f $(docker ps -q -a)
  docker system prune -f
}

function drmi() { docker rmi -f $(docker images -q); }

function drmall() {
  docker service rm $(docker service ls -q)
  docker rm -f $(docker ps -q -a);
  docker system prune -f
  docker rmi -f $(docker images -q);
  docker network rm $(docker network ls -q -f name=com.wedeploy)
  docker volume rm $(docker volume ls -qf dangling=true)
}
