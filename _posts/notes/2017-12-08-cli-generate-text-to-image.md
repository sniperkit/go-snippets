---
layout: note
title: "ImageMagick gen. text to image file"
tagline: ""
published: true
category : notes
tags : [note, cli, graphics]
---

Generate an image with my email address, from the command line :

    $ convert -background 'rgba(0,0,0,0)' -fill black \
        -font Source-Code-Pro-for-Powerline -pointsize 12 \
        label:"noreply@examble.com" \
        cadet.fabien_at_gmail.black.png

Likewise for my mobile phone number :

    $ convert -background 'rgba(0,0,0,0)' -fill black \
        -font Source-Code-Pro-for-Powerline -pointsize 12 \
        label:"+33 (0) 123 11 22 33" \
        cadet_fabien_mobile_phone_number.black.png

List available fonts :

    $ convert -list font | grep Font:

There seems to exist a possibly “better” alternative to ImageMagick,
probably a fork of it : __GraphicsMagick__.  But it behaves slightly differently.

    $ sudo apt-get install graphicsmagick

    $ gm convert -background 'rgba(0,0,0,0)' -fill white \
        -font Source-Code-Pro-for-Powerline -pointsize 12 \
        label:"+33 (0) 123 11 22 33" \
        cadet_fabien_mobile_phone_number.white.png
