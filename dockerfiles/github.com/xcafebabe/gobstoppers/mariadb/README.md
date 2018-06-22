# MariaDB 10 server on Alpine

## Starting the container

To run the lastest stable version of this docker image run

```
docker run --rm \
 -e MARIADB_RANDOM_ROOT_PASSWORD=yes \
 kafebob/rpi-mariadb
```

to expose the database to the external interface run

```
docker run --rm --name=mariadb-container \
-p 3306:3306 \
-e MARIADB_RANDOM_ROOT_PASSWORD=yes \
kafebob/rpi-mariadb
```

to assign a specific root password

```
docker run --rm --name=mariadb-container \
-p 3306:3306 \
-e MARIADB_ROOT_PASSWORD=umpalumpa \
-e MARIADB_REMOTE_ROOT=yes \
kafebob/rpi-mariadb
```

## Environment variables used in the container

### MARIADB_ROOT_PASSWORD
This variable defines the password for the root user in the database, se it with

	-e MARIADB_ROOT_PASSWORD=secretpassword

add quotes if there is spaces or other special character in the passwordd

	-e MARIADB_ROOT_PASSWORD='password with spaces'

### MARIADB_RANDOM_ROOT_PASSWORD
This variable generate a random password for the root user, add

	-e MARIADB_RANDOM_ROOT_PASSWORD=yes

the password can then be found by looking at the logoutput

	docker logs <container>

### MARIADB_ALLOW_EMPTY_PASSWORD
This allowes the root password to be blank, THIS IS A MAJOR SECURITY RISK, add

	-e MARIADB_ALLOW_EMPTY_PASSWORD=yes

### MARIADB_REMOTE_ROOT
Normal the root user can only use localhost to access the databases adding

	-e MARIADB_REMOTE_ROOT=yes

allows root access from any host

### MARIADB_DATABASE
creates a database with the defined name

	-e MARIADB_DATABASE=databasename

### MARIADB_USER
creates a user with password defined with MARIADB_PASSWORD and full access to the database defined by MARIADB_DATABASE

	-e MARIADB_USER=username

### MARIADB_PASSWORD
The apssword for the user defined by MARIADB_USER

	-e MARIADB_PASSWORD=donottell


### MARIADB_PASSWORD

## Volumes

The /data volume is defined containing

### /data/conf

Contains the configuration of mariadb (**my.cnf**)

### /data/db

Contains the database files

### /data/logs

Contains logs from mariadb
