#!/bin/bash

#this is constant and won't be accepted as a parameter
DB_IMAGE=postgres:9.5.9

#the parameters default values
db_user=dbadmin
db_pass=dbadmin
db_name=workshop_db
db_port=5432
container_name=postgres_db
force=false

#prints the help message
print_help(){
    echo "List of available parameters:"
    echo "--db_user             Default user that can be used to accessed the database. (default: $db_user)"
    echo "--db_pass             Password of the user used to access the database. (default: $db_pass) "
    echo "--db_name             Name of the database. (default:$db_name)"
    echo "--db_port             The port that will be exposed from the container. (default: $db_port)"
    echo "--container_name      The name that will be assigned to the container. (default: $container_name)"
    echo "--force -f            If passed, it will force the db removal."
    echo "--help -h             Prints this help message."
}

#function used to start the postgres container
start_container(){
    container_id=$(docker run --name $container_name -p $db_port:5432 -e \
        POSTGRES_USER=$db_user -e \
        POSTGRES_PASSWORD=$db_pass -e \
        POSTGRES_DB=$db_name -d $DB_IMAGE)

    echo "Postgres container started."
    echo "Container id:$container_id"
    echo "Db name     :$db_name"
    echo "DB user     :$db_user"
    echo "DB password :$db_pass"
    echo "Use 'docker log $container_name' to check the container log."
    exit 0
}

#read the parameters
while [ $# -gt 0 ]; do
  case "$1" in
    --db_user=*)
      db_user="${1#*=}"
      ;;
    --db_pass=*)
      db_pass="${1#*=}"
      ;;
    --db_name=*)
      db_name="${1#*=}"
      ;;
    --db_port=*)
      db_port="${1#*=}"
      ;;
    --container_name=*)
      container_name="${1#*=}"
      ;;
    -f|--force)
      force=true
      ;;
    -h|--help)
      print_help
      exit 1
      ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid parameters.*\n"
      printf "***************************\n"
      printf "Use -h or --help for the list of accepted parameters.\n"
      exit 1
  esac
  shift
done

#check if docker is installed
which docker

if [ $? -eq 0 ]
then
    docker --version | grep "Docker version"
    if [ $? -ne 0 ]
    then
        echo "This script will start the database as a docker container."
        echo "Please install docker."
        exit 1
    fi
else
    echo "This script will start the database as a docker container."
    echo "Please install docker."
    exit 1
fi

if [ ! "$(docker ps -a -q -f name=$container_name)" ]; then
    #the container does not exist so it can be started
    start_container
else
    if [ "$force" = true ]
    then
        removed_container=$(docker rm -f $container_name)
        echo "Container '$removed_container' removed."
        start_container
    else
        echo "Container '$container_name' already exists."
        read -r -p "Do you want to remove the container? [Y/n] " response
        #if the answer was y then continue with container removal and creation
        if [ "$response" = "Y" ]
        then

            removed_container=$(docker rm -f $container_name)
            echo "Container '$removed_container' removed."
            start_container
        else
            echo "Container creation aborted."
            exit 0
        fi
    fi
fi
