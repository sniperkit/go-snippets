# alpine-nginx

## Run Container
In order to run a container you need execute simple command
```
docker run -d --name alpine-nginx -p 80:80 -v /path/to/projects/:/var/www/html/ lykhouzov/alpine-nginx
```
Where path /path/to/projects/ is a local path to your projects

Generally i have following dir structure

```
/path/to/projects/
                project1/
                project2/
                ...
                projectN/
```


Also, you can add own config file for nginx.
```
docker run -d --name alpine-nginx -p 80:80 -v /home/vhosts/:/var/www/html/ -v ./php/example.conf:/etc/nginx/conf.d/www.conf lykhouzov/alpine-nginx
```

## Nginx configuration
php-fpm.conf file contains configuration which i use for my project.
It allow me setup new project very fast - i just need to add new directory in /path/to/projects/ and ... that's it. After that i can go to http://newproject.loc/

PS: i work on linux with installed dnsmask which configured to route all domains with .loc to 127.0.0.1
