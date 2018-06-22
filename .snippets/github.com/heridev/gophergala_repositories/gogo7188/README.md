## What's This?
My Gopher Gala Project.

This is minimal tool, that my Kindle Clippings on the Web.

http://kindleclippings.hexacosa.net/


## Team
solo project

* [Hideo Hattori(hhatto)](https://github.com/hhatto)


## Requirements
* [github.com/zenazn/goji](https://github.com/zenazn/goji)
* [github.com/lidashuang/goji_gzip](https://github.com/lidashuang/goji_gzip)
* [github.com/flosch/pongo2](https://github.com/flosch/pongo2)
* [github.com/hhatto/klip](https://github.com/hhatto/klip)


## Internals
```
    Local                                     on-premises server
 +-------------------+  BitTorrent Sync  +------------------------+
 |        Mac        |---------~|~-------|         Linux          |
 | [BitTorrent Sync] |                   |   [My Clippings.txt]   |
 +-------------------+                   |    |                   |
           |                             | +==|=================+ |
           |                             | |  |  gogo7188      <-------Internet[HTTP]
 +--------------------+                  | | (Go/goji/pongo2..) | |
 |       kindle       |                  | +====================+ |
 | [My Clippings.txt] |                  +------------------------+
 +--------------------+

```

