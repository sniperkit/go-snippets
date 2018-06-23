---
layout: page
title: todo, ideas, ...
tagline: Random brain dumps
published: false
---

Here lies some things I ought to write about, or not, ideas, thoughts, etc...
collected while coding or when my mind wanders between two shell commands.

## Symfony2

* [FOSRestBundle](https://github.com/FriendsOfSymfony/FOSRestBundle)

    - Out there: [REST APIs with Symfony2: The Right Way, 2012-08-02 by William Durand](http://williamdurand.fr/2012/08/02/rest-apis-with-symfony2-the-right-way/) ;

    - Basic setup ;

    - First steps ;

## JavaScript

* [fabic/jessme @ GitHub](https://github.com/fabic/jessme) : What was I trying to do with that thing, actually ?
* [Journey Through The JavaScript MVC Jungle, 2012-07-27 by Addy Osmani](http://coding.smashingmagazine.com/2012/07/27/journey-through-the-javascript-mvc-jungle/)
* [Learning JavaScript Design Patterns, a book by Addy Osmani © 2012](http://addyosmani.com/resources/essentialjsdesignpatterns/book/) : A great online book on design patterns in general, notably MVC/MVP/MVVM, with regards to JavaScript.
* [TodoMVC : Helping you select an MV* framework](http://addyosmani.github.com/todomvc/)

## PHP

* [ApiGen.php](http://apigen.org/) : _API documentation generator for PHP 5.3+_

    - ApiGen & [Symfony2](http://symfony.com/) / [Composer.phar](http://getcomposer.org/) : Your
      script for listing source directories & generating API doc for your app.
      along with dependencies.

## Linux / Gentoo / Debian

### Gentoo

##### Dude, from time to time you want to upgrade your Gentoo box :

    time \
        emerge --sync &&
        emerge --update --deep --with-bdeps=y --newuse --ask --keep-going world ; echo RETVAL=$?

    emerge --depclean --ask

    time revdep-rebuild -v

##### Linux kernel

    zcat /proc/config.gz > /usr/src/`date.....`
    genkernel ...

## Sybase ASE

* DB schema queries ;
* DB dump & restore ;

## Miscellaneous ramblings & wanderings

* [OrientDB](http://www.orientdb.org/), [ArangoDB](http://www.arangodb.org/) : Have a look at these graph / document databases, NoSQL things...

    - [doctrine/orientdb @ GitHub](https://github.com/doctrine/orientdb-odm) : _A set of tools to use and manage any OrientDB instance from PHP._
    - [triAGENS/ArangoDB-PHP @ GitHub](https://github.com/triAGENS/ArangoDB-PHP) : _A PHP client for ArangoDB._

* [spf13/spf13-vim](https://github.com/spf13/spf13-vim) : _Steve Francia's Vim Distribution_

