# <i>View Later</i> is a Chrome extension to keep links you want (maybe) to visit later



View Later is inspired by ReadLater extension (https://github.com/napsternxg/ReadLater).<br>
It allows you to pin the link you would like to visit later. <br>
ViewLater Chrome extension is written in Go with the usage of GopherJS. <br>

**Update 1:** Migrated to Bootstrap.

**Update 2:** Screenshots are resized in order to minimize storage size.   

<img width="564" alt="screen shot 2017-12-31 at 12 22 24 am" src="https://user-images.githubusercontent.com/3159236/34457933-dd6dc174-edc0-11e7-895b-a93982846b1b.png">

# How to install?
1. Download the sources.
2. Init govendor:
```
govendor init
govendor fetch github.com/fabioberger/chrome
govendor fetch honnef.co/go/js/dom
```
3. Build the extension:
```
govendor fetch honnef.co/go/js/dom
```
As result viewlater/ folder will be created. 

4. Go to chrome://extensions/ and click on "Load unpacked extension..." button.
Navigate to viewlater/ folder and click "Select" button. Extension will be added to you Chrome browser.

# How to use?
When on the tab you want to view later click on ViewLater extension button and press "Add".
Want to visit previously pinned link? Click on ViewLater extension button and then click on the link to visit.
Screenshot of the website is created for convenience.

# Problems with gopher.js and vendor
In case if you hit the issue https://github.com/gopherjs/gopherjs/issues/415
Please try this workaround from console (from the root of the project):
```
>cd vendor
>ln -s . src
>cd ..
```
