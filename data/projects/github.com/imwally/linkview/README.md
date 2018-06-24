# linkview

Like [urlview](https://github.com/sigpipe/urlview) but with more context for HTML documents.

That is, `linkview` parses both plaintext and valid HTML documents to
extract links. The advantage to this is being able to extract the text
from a link or in the case of an image, the alt or title
attributes. This gives more context to what the URL points to. Link
text will show in the menu instead of just the URL if given an HTML
document.

## Example

```
h: help   q: quit   (6 of 29)

https://eventing.coursera.org/redirectSigned/eyJrZXkiOiJlbWFpbC5saW5rLm9wZW4iLCJ2YWx1ZSI6

   Expert endorsed recommendations from a catalog of 2000+ courses -- little to no previo
   Python for Everybody
   Enroll Now
   Ruby on Rails
   Enroll Now
-> Algorithms
   Enroll Now
   Java Programming and Software Engineering Fundamentals
   Enroll Now
   Android App Development
   Enroll Now
   Full Stack Web and Multiplatform Mobile App Development
   Enroll Now
   Applied Data Science with Python
   Enroll Now
   Data Warehousing for Business Intelligence
   Enroll Now
   Cloud Computing
   Enroll Now
   MCS-DS 
   Apply Now
   FB
   Twitter
   LI
   iOS
   Android
   Learner Help Center |
   Email Settings |
   Unsubscribe
```

The first section displays some help text, the second displays a URL
preview of the currently selected link, and finally the menu of
links. Pressing the `return` key will open the currently selected link
in your default browser. Pressing the `tab` key will hide the menu and
display the full URL.

## How To Install

```
$ go get -u github.com/imwally/linkview
```

## How To Use

```
$ linkview /path/to/filename
```

Or pipe to.

```
$ cat /path/to/filename | linkview
```

## Help

```
h:               toggle help (press again to return to menu)
tab:             toggle full url
g:               go to top
G:               go to bottom
k / C-p / up:    move up
j / C-n / down:  move down
return / C-o:    open url
q / C-c:         quit
```

## But Why?

Many HTML emails are built using marketing platforms that modify URLs
for tracking and analytics purposes. It becomes difficult to figure
out what the URL points to. If you notice in the example above, the
URL to the Algorithms class is obscured by a redirect URL but
`linkview` displays the link text instead of only the URL.
