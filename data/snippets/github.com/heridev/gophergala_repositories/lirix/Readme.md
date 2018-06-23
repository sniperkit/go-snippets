Lirix
=====

![Logo](logo_100.png)

Lirix is a web-based weather app that uses data from OpenWeatherMap. It is written entirely in Go (and HTML for page templates) and does not use any third-party packages.

Issues
------

* Adding and removing cities has been (publicly) disabled due to database locking issue.
* All timezones are local to the computer (rather than the location).
* Search is implemented but not available because it doesn't actually do anything (see adding & removing cities issue).