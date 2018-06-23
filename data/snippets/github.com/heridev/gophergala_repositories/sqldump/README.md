# sqldump

A small tool for assisting database administration. My first 48 hours in Golang.

Unfortunately I do not have that much time, so it's just close to the backend with some basic html. 
Use your fantasy for best UX and choose high levels of abstraction and imagination for fancy output with latest js-technology. 

## prepare

    sudo mysqladmin --defaults-file=/etc/mysql/debian.cnf create galadb
    sudo mysql --defaults-file=/etc/mysql/debian.cnf -e "GRANT ALL PRIVILEGES  ON galadb.*  TO 'galagopher'@'localhost' IDENTIFIED BY 'mypassword'  WITH GRANT OPTION;"
    mysql -p"mypassword" -u galagopher galadb -e 'create table posts (title varchar(64) default null, start date default null);'
    mysql -p"mypassword" -u galagopher galadb -e 'insert into posts values("hello","2015-01-01");'
    mysql -p"mypassword" -u galagopher galadb -e 'insert into posts values("more","2015-01-03");'
    mysql -p"mypassword" -u galagopher galadb -e 'insert into posts values("end","2015-01-23");'
    mysql -p"mypassword" -u galagopher galadb -B -e 'select * from posts;'

## install

    export GOPATH=$PWD
    go get github.com/go-sql-driver/mysql
    go get github.com/gorilla/securecookie
    go get -u github.com/gophergala/sqldump

## run

    bin/sqldump


## usage

[http://localhost:8080](http://localhost:8080)

or more convenient but not safe

[http://localhost:8080/?user=galagopher&pass=mypassword&host=localhost&port=3306](http://localhost:8080/?user=galagopher&pass=mypassword&host=localhost&port=3306)


## caveats

- just basic protection against sql injection via database and table names

## perspectives

- choice for different database drivers
- insert and edit records

# License

MIT License
