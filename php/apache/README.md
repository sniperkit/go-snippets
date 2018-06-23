```
docker run -d -v ./:/var/www/ilsorriso \
    -v ./docker/httpd/vhost.conf:/etc/apache2/sites-available/000-default.conf \
    -v ./docker/httpd/php.ini:/etc/php5/apache2/php.ini \
    gianarb/php:apache
```
