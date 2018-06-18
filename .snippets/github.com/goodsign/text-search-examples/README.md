About
====================

This is an example of libtextcat+snowball+icu bundle usage for searching text in documents.

Installation
====================

* Install "github.com/goodsign/libtextcat"
* Install "github.com/goodsign/snowball"
* Install "github.com/goodsign/icu"
* go get github.com/goodsign/text-search-examples

After that you go to **$GOPATH/src/github.com/goodsign/text-search-examples** and

* go build
* ./text-search-examples

Detailed description
=====================

This example demonstrates how to search words in different languages in files with different language/encoding contents.

Searched words:

* supply 
* name 
* mein  (german)
* гора  (russian)

Searched files:

* eng_88591.txt (Language: english, encoding: ISO-8859-1)
* german_utf16le.txt (Language: german, encoding: UTF16LE with BOM)
* rus_koi8r (Language: russian, encoding: KOI8-R)

Example output:
----------------------

```
./text-search-examples
### FILE: input-documents/eng_88591.txt
    Detected encoding:
        ISO-8859-1
    Detected languages:
        english

    Found words:
    named.     (Original word: 'name' Stem: 'name')
    supplies.     (Original word: 'supply' Stem: 'suppli')
    name.     (Original word: 'name' Stem: 'name')

### FILE: input-documents/german_utf16le.txt
    Detected encoding:
        UTF-16LE
    Detected languages:
        german

    Found words:
    meine.     (Original word: 'mein' Stem: 'mein')
    meines.     (Original word: 'mein' Stem: 'mein')
    meine.     (Original word: 'mein' Stem: 'mein')
    mein.     (Original word: 'mein' Stem: 'mein')
    meine.     (Original word: 'mein' Stem: 'mein')
    meines.     (Original word: 'mein' Stem: 'mein')
    meinem.     (Original word: 'mein' Stem: 'mein')
    mein.     (Original word: 'mein' Stem: 'mein')
    meine.     (Original word: 'mein' Stem: 'mein')
    meiner.     (Original word: 'mein' Stem: 'mein')

### FILE: input-documents/rus_koi8r.txt
    Detected encoding:
        KOI8-R
    Detected languages:
        russian

    Found words:
    горами.     (Original word: 'гора' Stem: 'гор')
    гор.     (Original word: 'гора' Stem: 'гор')
```
