
     #    #  #    #  #####   ######
     #   #   #    #  #    #      #
     ####    #    #  #    #     #
     #  #    #    #  #####     #
     #   #   #    #  #   #    #
     #    #   ####   #    #  ######

_Kurz_ is yet another URL shortener/aliaser.

  * built as a Go data model API, for internal site use as well as public use
  * web front end separate from data logic
  * multiple aliasing strategies. provided:
    * identity
    * 32-bits hexa string, leading 0s omitted
    * manual selection
    * easy to extend, by registering additional AliasStrategy objects
  * supporting
    * vanity domains
    * relative paths (API only, no Web UI), for use within web sites
    * "slug"-type multipart aliasing for SEO
  * usage statistics
  * some tests available

It is currently available under the General Public License version 3 or later.


Using Kurz
==========

The main Kurz command show how to initialize the package and access statistics.

You can find usage examples in strategy_test.go:

- TestBaseAlias() shows how to generate aliases
- TestUseCounts() show how to query strategy statistics


Running tests
=============

Currently, only the strategy package has tests, and these need some setup, as
they touch the database:

- either redefine the credentials in strategy_test.go#initTestStorage()
- or create a database with privileges appropriate for them
- or send a PR submitting a better DB initialization mechanism for tests
