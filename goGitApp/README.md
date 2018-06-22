# goLangApp
for the gopher gala '15 - go code!

To see this code run live on your local machine,
follow the go installation instructions at <a href="http://golang.org">golang</a>
or similar.

With that installed, make a directory "..xyz" then cd xyz, place goDB.go in it.

If you're on windows and go is on your PATH already:
<code>set GOPATH="./xyz"</code>

Next,simply <code>go run goDB.go</code>, provided mysql is running on your
system already. You can use WAMPServer for that part if needed, too, since
it makes it easy to set up the root password and port using the system tray
icon for wampserver.  

Example sql commands to set up your table might be:
create database test;
use test;
create table links(id int(9) NOT NULL PRIMARY KEY AUTO_INCREMENT, name varchar(80));
insert into links(id,name) values(0,'steve');

From a mysql:// prompt you can use "explain links" to see the structure of your
mysql table.