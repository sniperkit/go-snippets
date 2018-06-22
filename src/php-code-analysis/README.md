# Dockerfile for PHP code analysis

![](https://img.shields.io/docker/pulls/qmu1/php-code-analysis.svg)
![](https://img.shields.io/docker/build/qmu1/php-code-analysis.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/php-code-analysis.svg)](https://microbadger.com/images/qmu1/php-code-analysis)

An alipine-based dockerfile which includes code analysis tools for PHP.

* [phpcs (PHP_CodeSniffer)](https://github.com/squizlabs/PHP_CodeSniffer)
* [phpmd (PHP Mess Detector)](https://github.com/phpmd/phpmd)
* [phpcpd (PHP Copy/Paste Detector)](https://github.com/sebastianbergmann/phpcpd)
* [php-cs-fixer (PHP Code Standars Fixer)](https://github.com/FriendsOfPHP/PHP-CS-Fixer)
* [php-short-array-syntax-converter](https://github.com/thomasbachem/php-short-array-syntax-converter)

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/php-code-analysis/Dockerfile)

## Usage

```
docker run -it --rm -v $(pwd):/src qmu1/php-code-analysis:latest /bin/sh -c "\
    phpcs --version;
    phpmd --version;
    phpcpd --version;
    php-cs-fixer --version;
"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/php-code-analysis/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
