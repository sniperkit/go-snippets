A golang-based digital signage server.
===

The goal of this project is a simple, easy-to-implement digital signage server that can be utilized with nothing more than a web browser on the client side.

Example
=======

Head to http://agocs.org:3030/ to see this thing running live!

Check out http://agocs.org:3030/config.html to play with the configuration. Changes you make affect everybody, so be nice. If someone does something objectionable or you need help configuring it, email me. chris@agocs.org.

The current known-good configuration is:

	http://sshchicago.org | http://pumpingstationone.org
	https://www.google.com/calendar/embed?height=600&wkst=1&bgcolor=%23FFFFFF&src=ssh%40sshchicago.org&color=%23B1440E&ctz=America%2FChicago | https://www.google.com/calendar/embed?mode=AGENDA&height=600&wkst=1&bgcolor=%23FFFFFF&src=hhlp4gcgvdmifq5lcbk7e27om4%40group.calendar.google.com&color=%23A32929&ctz=America%2FChicago
	=
	http://forecast.io/embed/#lat=41.8915&lon=-87.6146&name=Chicago&color=#00aaff&font=Georgia&units=us | http://www.cityofchicago.org/city/en/depts/mayor/iframe/plow_tracker.html
	=
	http://www.theonion.com | http://edition.cnn.com/?hpt=header_edition-picker
	http://www.foxnews.com | http://www.msnbc.com/
	=
	http://comeonandsl.am/


Usage
===
The server runs, by default, on port 3030 of your host.

The configuration page lives at <server root>/config.html and allows you to enter a series of URLS in a layout language defined below. There may be a pause of up to a minute long before the server begins providing these URLs.

The client lives at the server root, and is simply a small amount of javascript that loads your pages to display in iframes. The border and scrollbar has been removed - pages that do not require scrolling are currently required.

Page Layout DSL
===============

The DSL we've created has three special characters:

- the ` | ` (space pipe space), which is used to separate columns
- the `\n` (newline character), which is used to separate rows
- the `=` (equals), which is used to separate pages

A display has one or more pages, each of which has one or more rows with one or more columns of urls to load. 

Example: Let's say you want two pages. The first page will show a list of events and weather side-by-side, and traffic on a new line. The second page will show only your website. Your configuration would look like this:

	http://list-of-events.com/yourOrg | http://forecast.io/#/f/your_location
	http://traffic.com/your_location
	=
	http://yourOrg.org

Parsing this DSL is incredibly non-fault-tolerant right now, so don't mess up.