---
layout: post
title: "Building PHP (from its Git repository)"
description: ""
category: diary
tags: [PHP, make, old]
---

_Saturday, windy and cold day, went out to get some coffee on a café's terrace... and went back chilled and chilled out_


## Git clone it!

* <http://www.php.net/git.php> ;
* https://github.com/php/php-src

        (12:44:07) ○ [fabi@winterfell] ~/dev  git clone https://github.com/php/php-src.git
        Cloning into 'php-src'...
        remote: Finding bitmap roots...
        remote: Reusing existing pack: 526255, done.
        remote: Counting objects: 5557, done.
        remote: Compressing objects: 100% (2881/2881), done.
        remote: Total 531812 (delta 3837), reused 4127 (delta 2672)
        Receiving objects: 100% (531812/531812), 148.56 MiB | 1.52 MiB/s, done.
        Resolving deltas: 100% (410336/410336), done.
        (12:48:13) ○ [fabi@winterfell] ~/dev

## Checkout a stable branch :

        (12:57:05) ± [fabi@winterfell] ~/dev/php-src (master ✓)  git checkout -b php-5.4.22 tags/php-5.4.22
        Switched to a new branch 'php-5.4.22'
        (12:57:24) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)


## About autoconf, automake & libtool versions

See <http://www.php.net/git.php>.

        (13:00:23) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)  libtool --version
        libtool (GNU libtool) 2.4.2

        (13:01:07) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)  autoconf --version
        autoconf (GNU Autoconf) 2.69

        (13:01:16) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)  automake --version
        automake (GNU automake) 1.13.4

## ./buildconf

        (13:05:43) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)  ./buildconf --force
        Forcing buildconf
        Removing configure caches
        buildconf: checking installation...
        buildconf: autoconf version 2.69 (ok)
        rebuilding aclocal.m4
        rebuilding configure
        rebuilding main/php_config.h.in
        (13:05:53) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)

        (13:06:24) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✓)  ./configure --help |less
        ...
        ...

## "auto" target dir. prefix :

        (13:27:57) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)  git branch
          master
        * php-5.4.22

        (13:28:03) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)  branch=`git symbolic-ref HEAD` && branch=${branch##refs/heads/} && prefix=/opt/$branch ; echo "branch=$branch"; echo prefix=$prefix;
        branch=php-5.4.22
        prefix=/opt/php-5.4.22

        (13:28:05) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)

## ./configure --prefix=...


* `--enable-zend-multibyte` & `--enable-sqlite-utf8` : Not for >= 5.4.\*+

        # prefix=/opt/$(basename `pwd`)

        time \
        ./configure --prefix=${prefix:?'HEY! Forgot to provide an installation --prefix ??'} \
            --with-apxs2=/usr/sbin/apxs2 \
            --enable-zend-multibyte \
            --enable-mod-charset \
            --enable-mbstring \
            --enable-intl \
            --enable-shmop --enable-sysvmsg --enable-sysvsem --enable-sysvshm --enable-pcntl \
            --enable-calendar --enable-exif --enable-ftp \
            --with-gd --with-freetype-dir=/usr --enable-gd-native-ttf \
            --with-jpeg-dir=/usr --with-png-dir=/usr \
            --enable-soap --enable-sockets \
            --enable-zip --with-zlib --with-bz2 \
            --with-readline \
            --with-openssl --with-curl \
            --enable-wddx \
            --with-xsl \
            --enable-dba=shared --enable-sqlite-utf8 \
            --with-mysql=mysqlnd --with-mysqli=mysqlnd --with-pdo-mysql=mysqlnd \
            --with-mysql-sock=/var/run/mysqld/mysqld.sock \
            --with-sybase-ct=/opt/sybase/OCS-15_0 \
            --with-tidy \
            2>&1 | tee configure.log.txt

## make && ~ install

### Make

        $ time make -j3

### Install

Using `-k` (`--keep-going`) for it won't have write access to whence Apache modules lies :

        $ make install -k 2>&1 | tee install_`date -I`.log.txt

#### Output in error is along these lines:

        Installing PHP SAPI module:       apache2handler
        /usr/lib64/apache2/build/instdso.sh SH_LIBTOOL='/usr/bin/libtool' libphp5.la /usr/lib64/apache2/modules
        /usr/bin/libtool --mode=install cp libphp5.la /usr/lib64/apache2/modules/
        libtool: install: cp .libs/libphp5.so /usr/lib64/apache2/modules/libphp5.so
        cp: cannot create regular file '/usr/lib64/apache2/modules/libphp5.so': Permission denied
        apxs:Error: Command failed with rc=65536
        .
        make: *** [install-sapi] Error 1

#### Copying libphp5.so manually to somewhat of an appropriate place, at least for me :

        ~/dev/php-src (php-5.4.22 ✗)  cp -iv .libs/libphp5.so $prefix/lib/
        ‘.libs/libphp5.so’ -> ‘/opt/php-5.4.22/lib/libphp5.so’

## php.ini files & further setup

        (14:18:15) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)  cp -iv php.ini-* $prefix/lib/
        ‘php.ini-development’ -> ‘/opt/php-5.4.22/lib/php.ini-development’
        ‘php.ini-production’ -> ‘/opt/php-5.4.22/lib/php.ini-production’

        (14:19:31) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)  pushd $prefix/lib
        /opt/php-5.4.22/lib ~/dev/php-src

        (14:20:11) ○ [fabi@winterfell] /opt/php-5.4.22/lib  ln -sv php.ini-development php.ini
        ‘php.ini’ -> ‘php.ini-development’

        (14:21:18) ○ [fabi@winterfell] 0 /opt/php-5.4.22/lib  ll
        total 35168
        drwxr-sr-x  3 fabi wheel     4096 Nov 23 14:20 .
        drwxr-sr-x  7 fabi wheel     4096 Nov 23 13:54 ..
        -rwxr-xr-x  1 fabi wheel 35858376 Nov 23 13:58 libphp5.so
        drwxr-sr-x 15 fabi wheel     4096 Nov 23 13:54 php
        lrwxrwxrwx  1 fabi wheel       19 Nov 23 14:20 php.ini -> php.ini-development
        -rw-r--r--  1 fabi wheel    66920 Nov 23 14:18 php.ini-development
        -rw-r--r--  1 fabi wheel    66952 Nov 23 14:18 php.ini-production
        (14:21:28) ○ [fabi@winterfell] 0 /opt/php-5.4.22/lib

## pecl install apc

        (14:45:52) ± [fabi@winterfell] 1 ~/dev/php-src (php-5.4.22 ✗)  pecl install apc
        downloading APC-3.1.13.tgz ...
        Starting to download APC-3.1.13.tgz (171,591 bytes)
        .....................................done: 171,591 bytes
        55 source files, building
        WARNING: php_bin /opt/php-5.4.22/bin/php appears to have a suffix -5.4.22/bin/php, but config variable php_suffix does not match
        running: phpize
        Configuring for:
        PHP Api Version:         20100412
        Zend Module Api No:      20100525
        Zend Extension Api No:   220100525
        Enable internal debugging in APC [no] :
        Enable per request file info about files used from the APC cache [no] :
        Enable spin locks (EXPERIMENTAL) [no] :
        Enable memory protection (EXPERIMENTAL) [no] :
        Enable pthread mutexes (default) [no] :
        Enable pthread read/write locks (EXPERIMENTAL) [yes] :
        building in /tmp/pear/temp/pear-build-fabiFcLiBG/APC-3.1.13
        running: /tmp/pear/temp/APC/configure --enable-apc-debug=no --enable-apc-filehits=no --enable-apc-spinlocks=no --enable-apc-memprotect=no --enable-apc-pthreadmutex=no --enable-apc-pthreadrwlocks=yes
        ...
        ...
        Build process completed successfully
        Installing '/opt/php-5.4.22/lib/php/extensions/no-debug-non-zts-20100525/apc.so'
        Installing '/opt/php-5.4.22/include/php/ext/apc/apc_serializer.h'
        install ok: channel://pecl.php.net/APC-3.1.13
        configuration option "php_ini" is not set to php.ini location
        You should add "extension=apc.so" to php.ini

## pecl install xdebug

        (14:46:19) ± [fabi@winterfell] 1 ~/dev/php-src (php-5.4.22 ✗)   pecl install xdebug
        downloading xdebug-2.2.3.tgz ...
        Starting to download xdebug-2.2.3.tgz (250,543 bytes)
        .....................................................done: 250,543 bytes
        66 source files, building
        WARNING: php_bin /opt/php-5.4.22/bin/php appears to have a suffix -5.4.22/bin/php, but config variable php_suffix does not match
        running: phpize
        Configuring for:
        PHP Api Version:         20100412
        Zend Module Api No:      20100525
        Zend Extension Api No:   220100525
        building in /tmp/pear/temp/pear-build-fabiA9Zp6e/xdebug-2.2.3
        running: /tmp/pear/temp/xdebug/configure
        ...
        ...
        Build process completed successfully
        Installing '/opt/php-5.4.22/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so'
        install ok: channel://pecl.php.net/xdebug-2.2.3
        configuration option "php_ini" is not set to php.ini location
        You should add "zend_extension=xdebug.so" to php.ini

## php.ini : Some dev. configurations


        zend_extension=/opt/php-5.4.22/lib/php/extensions/no-debug-non-zts-20100525/xdebug.

        extension=apc.so

        xdebug.remote_enable = On

        ; For Composer.phar & Symfony2 notably :
        detect_unicode = Off

        ; Ah! legacy code...
        short_open_tag = On

        ; Huge PDFs generation may take some time with DOMPDF :
        max_execution_time = 600

        ; For ApiGen.php which tends to get hungry at times...
        memory_limit = 6G

        ; Legacy codes may need $_ENV :
        variables_order = "EGPCS"

        ; default_charset : Can't tell which is the default, for which PHP version :/
        ; PHP 5.x.y ?
        ;default_charset = "iso-8859-1"

        ; PHP 5.4.x ?
        ;default_charset = "UTF-8"

        ; Setting it "" so that no Content-Type HTTP header gets sent per default :
        ;default_charset = ""

        date.timezone = Europe/Paris

        error_reporting = E_ALL
        display_errors = On
        display_startup_errors = On
        log_errors = On
        html_errors = Off

        post_max_size = 10M
        upload_max_filesize = 8M

## make clean/distclean : for further compilations...

    (15:17:22) ± [fabi@winterfell] 0 ~/dev/php-src (php-5.4.22 ✗)  make distclean
    find . -name \*.gcno -o -name \*.gcda | xargs rm -f
    find . -name \*.lo -o -name \*.o | xargs rm -f
    find . -name \*.la -o -name \*.a | xargs rm -f
    find . -name \*.so | xargs rm -f
    find . -name .libs -a -type d|xargs rm -rf
    rm -f libphp5.la sapi/cli/php libphp5.la modules/* libs/*
    rm -f Makefile config.cache config.log config.status Makefile.objects Makefile.fragments libtool main/php_config.h stamp-h sapi/apache/libphp5.module buildmk.stamp Zend/zend_dtrace_gen.h Zend/zend_dtrace_gen.h.bak
    /bin/grep -E define'.*include/php' /home/fabi/dev/php-src/configure | /bin/sed 's/.*>//'|xargs rm -f

## PHP v5.5.6 : let's give it a try!

        (15:19:46) ± [fabi@winterfell] ~/dev/php-src (php-5.4.22 ✗)  git checkout -b php-5.5.6 tags/php-5.5.6
        M   ext/sybase_ct/php_sybase_ct.h
        Switched to a new branch 'php-5.5.6'

        prefix=/opt/php...
        ./buildconf --force
        ./configure ...
        make...

        # APC didn't build :
        pecl install apc

        # Xdebug was fine:
        pecl install xdebug

        vimdiff /opt/php-5.4.22/lib/php.ini /opt/php-5.5.6/lib/php.ini

## Hummm... -_-

* <http://php.net/manual/en/migration53.php>
* <http://php.net/manual/en/migration54.php>
* <http://php.net/manual/en/migration55.php>
