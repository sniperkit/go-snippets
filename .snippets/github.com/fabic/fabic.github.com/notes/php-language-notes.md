---
#layout: default
layout: page
title: "PHP: personal notes"
tagline: "Some personal collected notes about programming in PHP."
category : notes
tags : [draft]
published: true
---

### Packages and stuff of general interest

* [symfony/property-access](https://packagist.org/packages/symfony/property-access)
  [doc](https://symfony.com/doc/current/components/property_access.html)


## Composer, packages

* <https://getcomposer.org/doc/04-schema.md#repositories>
  * <https://getcomposer.org/doc/05-repositories.md#path>
* <https://getcomposer.org/doc/05-repositories.md>
* [2017: Developing composer packages locally, johannespichler.com](https://johannespichler.com/developing-composer-packages-locally/)
* [2017: How To Make a PHP Package (Mahmoud Zalt, Medium)](https://medium.com/@Mahmoud_Zalt/how-to-make-a-php-package-5d3bb6b00df0)
  > Now I will share some standards I consider while developing new package, that you should consider too:
  > - Use __Semantic versioning__. And try not to introduce
  >   __backwards compatibility (BC) breaks__.
  > - Write __unit tests__. Aim for at least 60% __coverage__ in the first version.
  > - Try not introduce __dependencies__ on other packages.
  > - Conform to PSR-1 and PSR-2 as the coding style guide. (to apply the conventions on an existing package install PHP Code Sniffer.)
  > - Conform to __PSR-4 for autoloading__.
  > - DocBlock your code.
  > - Have an extensive README. And keep it updated.
  > - Default to the use of the Solar vocabulary for method names.
  > - Keep ‘Public’ functions names short (maximum of two words).
  > - Compose functionality through dependency injection, rather than through inheritance and base classes.
  > - Try to limit the levels of indentation per function (one per function is great).
  > - Use Continuous Integration tools, to automatically check coding standards and run tests.
  > - Prefer explicit configuration over implicit convention.
  > - Exclude non-essential files in .gitattributes.
  > - It’s always better to have a ‘Changelog’.
  > - Do never ever have duplicated code.
  > __Mahmoud Zalt__

__EOF__
