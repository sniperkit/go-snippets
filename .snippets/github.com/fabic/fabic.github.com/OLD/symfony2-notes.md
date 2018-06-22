---
layout: page
title: "Symfony2 notes"
description: "Miscellaneous Symfony2 notes"
published: false
---

_Just some notes & links laid down while discovering Symfony2._

## Miscellaneous links

* **<http://www.symfony2cheatsheet.com/>** ;
* <http://symfony.com/blog/new-in-symfony-2-2-new-propertyaccess-component>
* <http://sonata-project.org/bundles/easy-extends/master/doc/reference/why.html>

## Form component

* [GH:symfony/Form](https://github.com/symfony/Form) ;
* [Book:Forms](http://symfony.com/doc/current/book/forms.html) ;
* [Cookbook:How to Dynamically Modify Forms Using Form Events](http://symfony.com/doc/current/cookbook/form/dynamic_form_modification.html) ;
* [Cookbook:How to use Data Transformers](http://symfony.com/doc/current/cookbook/form/data_transformers.html)
* [Reference:Form Field Type Reference](http://symfony.com/doc/current/reference/forms/types.html) ;
* [Reference:Validation Constraints Reference](http://symfony.com/doc/current/reference/constraints.html) ;
* [Reference:Twig Template Function and Variable Reference](http://symfony.com/doc/current/reference/forms/twig_reference.html) ;

### Blog posts, articles, lecture slides

* **[2012-03-06 : Symfony2 Form Architecture, by Bernhard Schussek](http://webmozarts.com/2012/03/06/symfony2-form-architecture/)** ;
* [2012-11-23 : 3 Steps to Symfony2 Form Mastery, by B. Schussek // Speaker Deck](https://speakerdeck.com/bschussek/3-steps-to-symfony2-form-mastery) ;
* [2012-09-28 : Symfony2 Form Tricks, by B. Schussek // Speaker Deck](https://speakerdeck.com/bschussek/symfony2-form-tricks) ;
* [2012-03-13 : Small Talk about the Symfony2 Form Component, by Thomas Rabaix](http://www.slideshare.net/th0masr/small-prformcomponent) ;
* [2011-03-03 : Leveraging Symfony2 Forms // Speaker Deck, by B. Schussek](https://speakerdeck.com/bschussek/leveraging-symfony2-forms) ;

### Class diagrams (B. Schussek)

#### Form , FormBuilder, FormView

![Class diagram, by B. Schussek][100]

#### Type guessing

![Type guessing class diagram, by B. Schussek][101]

### Form types hierarchy (B. Shussek)

![Form types hierarchy][102]

### Howto...

    [ ] Handle association-class relationships.

## Twig

* [GH:fabpot/Twig](https://github.com/fabpot/Twig) ;
* [twig.sensiolabs.org/documentation](http://twig.sensiolabs.org/documentation) ;
* [Reference:Twig Extensions (forms, filters, tags, etc) Reference](http://symfony.com/doc/current/reference/twig_reference.html) ;

#### How to...

    [ ] Twig & having DB-backed templates.

## Dependency Injection (DI)

DI, DIC

fabpot's articles



## Doctrine 2

* [doctrine/doctrine2](https://github.com/doctrine/doctrine2) ;
* [doctrine/dbal](https://github.com/doctrine/dbal) ;
* [doctrine-project.org/projects/common.html](http://www.doctrine-project.org/projects/common.html) ;
* [doctrine-project.org/projects/orm.html](http://www.doctrine-project.org/projects/orm.html) ;
* [doctrine-project.org/projects/dbal.html](http://www.doctrine-project.org/projects/dbal.html) ;

## Sonata \*


## Symfony CMF

* [cmf.symfony.com](http://cmf.symfony.com/) (See source [GH:/symfony-cmf/symfony-cmf-website/](https://github.com/symfony-cmf/symfony-cmf-website/)) ;
* [symfony.com/doc/master/cmf/](http://symfony.com/doc/master/cmf/) ([documentation source (GH)](https://github.com/symfony-cmf/symfony-cmf-docs));

* <http://symfony.com/doc/master/cmf/getting-started/installing-symfony-cmf.html>
* <http://symfony.com/doc/master/cmf/tutorials/installing-cmf-core.html>
* <http://symfony.com/doc/master/cmf/tutorials/installing-configuring-doctrine-phpcr-odm.html>
* <http://symfony.com/doc/master/cmf/bundles/phpcr-odm.html>

* <http://symfony.com/doc/master/cmf/bundles/create.html>

* <http://decoupledcms.org/>

### ...

    $ ./console doctrine:database:create
    $ ./console doctrine:phpcr:init:dbal
    $ ./console doctrine:phpcr:register-system-node-types
    $ ./console doctrine:phpcr:fixtures:load

### ...

    ~/dev/symfony-cmf-se (master ✗)  ./console doctrine:phpcr:mapping:info
    Found 17 documents mapped in document manager:
    [OK]   Doctrine\ODM\PHPCR\Document\Folder
    [OK]   Doctrine\ODM\PHPCR\Document\File
    [OK]   Doctrine\ODM\PHPCR\Document\Generic
    [OK]   Doctrine\ODM\PHPCR\Document\Resource
    [OK]   Doctrine\ODM\PHPCR\Document\Image
    [OK]   Symfony\Cmf\Bundle\MenuBundle\Document\MenuNode
    [OK]   Symfony\Cmf\Bundle\MenuBundle\Document\MultilangMenuNode
    [OK]   Symfony\Cmf\Bundle\ContentBundle\Document\StaticContent
    [OK]   Symfony\Cmf\Bundle\ContentBundle\Document\MultilangStaticContent
    [OK]   Symfony\Cmf\Bundle\SimpleCmsBundle\Document\Page
    [OK]   Symfony\Cmf\Bundle\SimpleCmsBundle\Document\MultilangRedirectRoute
    [OK]   Symfony\Cmf\Bundle\SimpleCmsBundle\Document\MultilangPage
    [OK]   Symfony\Cmf\Bundle\SimpleCmsBundle\Document\MultilangRoute
    [OK]   Symfony\Cmf\Bundle\CreateBundle\Document\Image
    [OK]   Symfony\Cmf\Bundle\RoutingExtraBundle\Document\Route
    [OK]   Symfony\Cmf\Bundle\RoutingExtraBundle\Document\RedirectRoute
    [OK]   Symfony\Component\Routing\Route


    ~/dev/symfony-cmf-se (master ✗)  ./console doctrine:phpcr:workspace:list
    The following 1 workspaces are available:
    default


    ~/dev/symfony-cmf-se (master ✗)  ./console router:debug
    [router] Current routes
    Name                                 Method Scheme Host Path
    _assetic_051ded3                     ANY    ANY    ANY  /css/051ded3.css
    _assetic_051ded3_0                   ANY    ANY    ANY  /css/051ded3_tags_1.css
    _assetic_051ded3_1                   ANY    ANY    ANY  /css/051ded3_create-ui_2.css
    _assetic_051ded3_2                   ANY    ANY    ANY  /css/051ded3_font-awesome_3.css
    _assetic_051ded3_3                   ANY    ANY    ANY  /css/051ded3_midgardnotif_4.css
    _assetic_051ded3_4                   ANY    ANY    ANY  /css/051ded3_halloCmfStyle_5.css
    _assetic_051ded3_5                   ANY    ANY    ANY  /css/051ded3_createStyle_6.css
    _assetic_051ded3_6                   ANY    ANY    ANY  /css/051ded3_overlay_7.css
    _assetic_8dae729                     ANY    ANY    ANY  /js/hallo-extra.js
    _assetic_8dae729_0                   ANY    ANY    ANY  /js/hallo-extra_init-create-hallo_1.js
    _assetic_8dae729_1                   ANY    ANY    ANY  /js/hallo-extra_hallo-min_2.js
    _assetic_74598c3                     ANY    ANY    ANY  /js/create.js
    _assetic_74598c3_0                   ANY    ANY    ANY  /js/create_jquery-1.7.1.min_1.js
    _assetic_74598c3_1                   ANY    ANY    ANY  /js/create_jquery-ui-1.8.18.custom.min_2.js
    _assetic_74598c3_2                   ANY    ANY    ANY  /js/create_underscore-min_3.js
    _assetic_74598c3_3                   ANY    ANY    ANY  /js/create_backbone-min_4.js
    _assetic_74598c3_4                   ANY    ANY    ANY  /js/create_rangy-core-1.2.3_5.js
    _assetic_74598c3_5                   ANY    ANY    ANY  /js/create_vie-min_6.js
    _assetic_74598c3_6                   ANY    ANY    ANY  /js/create_jquery.rdfquery.min_7.js
    _assetic_74598c3_7                   ANY    ANY    ANY  /js/create_jquery.tagsinput.min_8.js
    _assetic_74598c3_8                   ANY    ANY    ANY  /js/create_annotate-min_9.js
    _assetic_74598c3_9                   ANY    ANY    ANY  /js/create_create-min_10.js
    _wdt                                 ANY    ANY    ANY  /_wdt/{token}
    _profiler_home                       ANY    ANY    ANY  /_profiler/
    _profiler_search                     ANY    ANY    ANY  /_profiler/search
    _profiler_search_bar                 ANY    ANY    ANY  /_profiler/search_bar
    _profiler_purge                      ANY    ANY    ANY  /_profiler/purge
    _profiler_info                       ANY    ANY    ANY  /_profiler/info/{about}
    _profiler_import                     ANY    ANY    ANY  /_profiler/import
    _profiler_export                     ANY    ANY    ANY  /_profiler/export/{token}.txt
    _profiler_phpinfo                    ANY    ANY    ANY  /_profiler/phpinfo
    _profiler_search_results             ANY    ANY    ANY  /_profiler/{token}/search/results
    _profiler                            ANY    ANY    ANY  /_profiler/{token}
    _profiler_router                     ANY    ANY    ANY  /_profiler/{token}/router
    _profiler_exception                  ANY    ANY    ANY  /_profiler/{token}/exception
    _profiler_exception_css              ANY    ANY    ANY  /_profiler/{token}/exception.css
    _configurator_home                   ANY    ANY    ANY  /_configurator/
    _configurator_step                   ANY    ANY    ANY  /_configurator/step/{index}
    _configurator_final                  ANY    ANY    ANY  /_configurator/final
    home_redirect                        ANY    ANY    ANY  /
    static                               ANY    ANY    ANY  /{_locale}/static
    symfony_cmf_create_put_document      PUT    ANY    ANY  /{_locale}/symfony-cmf/create/document/{subject}
    symfony_cmf_create_post_document     POST   ANY    ANY  /{_locale}/symfony-cmf/create/document/{subject}
    symfony_cmf_create_put_document_base ANY    ANY    ANY  /{_locale}/symfony-cmf/create/document
    symfony_cmf_create_image_upload      POST   ANY    ANY  /symfony-cmf/create/image/upload
    symfony_cmf_create_image_search      GET    ANY    ANY  /symfony-cmf/create/image/search/
    symfony_cmf_create_image_related     GET    ANY    ANY  /symfony-cmf/create/image/related/
    symfony_cmf_create_image_list        GET    ANY    ANY  /symfony-cmf/create/image/list/
    symfony_cmf_create_image_display     GET    ANY    ANY  /symfony-cmf/create/image/{name}


## TODOs

    [ ] Have a simple <frameset> base thing for quick browing of these links.
    [ ] On ApiGen-erating the whole API documentation from Composer.phar's autoloader ;
    [ ] Annotations howto ?
    [ ] REST...?
    [ ] <https://read-the-docs.readthedocs.org/en/latest/getting_started.html> ([GH:/rtfd/readthedocs.org](https://github.com/rtfd/readthedocs.org))

### Other pending questions, howto...




[100]: http://webmozarts.com/wp-content/uploads/2012/03/class-diagram2.png
[101]: http://webmozarts.com/wp-content/uploads/2012/03/form-factory.png
[102]: http://webmozarts.com/wp-content/uploads/2012/03/types.png
