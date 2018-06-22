---
layout: post
title: "GH:antialize/wkhtmltopdf"
description: "On building antialize/wkhtmltopdf, a really nice command line tool for producing PDFs from web pages :-)"
category:
tags: [PDF, make, old]
---

_On building [wkhtmltopdf](http://code.google.com/p/wkhtmltopdf/), a great command line tool for rendering PDFs from webpages, with abilities for concatenation and generation of table of contents._

* <http://code.google.com/p/wkhtmltopdf/> : Homepage ;
* <http://code.google.com/p/wkhtmltopdf/w/list> : Wiki ;
* [GH:/antialize/wkhtmltopdf](https://github.com/antialize/wkhtmltopdf) : wkhtmltopdf code ;
* [Gitorious:/~antialize/qt/antializes-qt](http://gitorious.org/~antialize/qt/antializes-qt) : Customized QT (4.8.4 as of February 10th 2013, last commit: [08c134d](http://qt.gitorious.org/~antialize/qt/antializes-qt/commit/08c134d0b8a541cbb7ab6396baff1f07fe437d98), optional).

Other tools out there :

* [pdftk, the PDF toolkit, by pdflabs.com](http://www.pdflabs.com/tools/pdftk-the-pdf-toolkit/) ;
* [dompdf](http://code.google.com/p/dompdf/) ([GH:/dompdf/dompdf](https://github.com/dompdf/dompdf)).

## Installation

See <http://code.google.com/p/wkhtmltopdf/wiki/compilation>.

### Building antialize-qt (a.k.a wkhtmltopdf-qt) :

    $ git clone git://gitorious.org/~antialize/qt/antializes-qt.git wkhtmltopdf-qt

    $ cd wkhtmltopdf-qt/
    $ git checkout -b qt484 origin/4.8.4

    # See output below.
    $ QTDIR=. ./bin/syncqt

    ### CONFIGURE
    # Humm... had to remove -noixes parameter (unknown from ./configure) :
    #
    # See ./configure ... tail of output below.
    #
    $ time \
        ./configure -prefix ../wkqt -release -static -fast -exceptions \
            -no-accessibility -no-stl -no-sql-ibase -no-sql-mysql -no-sql-odbc \
            -no-sql-psql -no-sql-sqlite -no-sql-sqlite2 -no-qt3support -xmlpatterns \
            -no-phonon -no-phonon-backend -webkit -no-scripttools -no-mmx -no-3dnow \
            -no-sse -no-sse2 -qt-zlib -qt-libtiff -qt-libpng -qt-libmng -qt-libjpeg \
            -graphicssystem raster -opensource \
            -nomake tools -nomake examples -nomake demos -nomake docs -nomake translations \
            -no-opengl -no-dbus -no-multimedia -openssl -no-declarative -largefile -rpath \
            -no-nis -no-cups -no-iconv -no-pch -no-gtkstyle -no-nas-sound -no-sm -no-xshape \
            -no-xinerama -no-xcursor -no-xrandr -no-mitshm -no-xinput -no-xkb \
            -no-glib -no-openvg -no-xsync -no-audio-backend \
             -confirm-license        2>&1 | tee conflog.txt

    ## Output :
            This target is using the GNU C++ compiler (linux-g++-64).

            Recent versions of this compiler automatically include code for
            exceptions, which increase both the size of the Qt libraries and
            the amount of memory taken by your applications.

            You may choose to re-run configure with the -no-exceptions
            option to compile Qt without exceptions. This is completely binary
            compatible, and existing applications will continue to work.


    Build type:    linux-g++-64
    Architecture:  x86_64
    Platform notes:

                - Also available for Linux: linux-kcc linux-icc linux-cxx

    Build .................. libs
    Configuration ..........  release static largefile avx  minimal-config small-config medium-config large-config full-config reduce_exports ipv6 clock-gettime clock-monotonic mremap getaddrinfo ipv6ifname getifaddrs inotify jpeg mng png gif tiff system-freetype zlib openssl xrender fontconfig concurrent xmlpatterns svg script release
    Debug .................. no
    Qt 3 compatibility ..... no
    QtDBus module .......... no
    QtConcurrent code ...... yes
    QtGui module ........... yes
    QtScript module ........ yes
    QtScriptTools module ... no
    QtXmlPatterns module ... yes
    Phonon module .......... no
    Multimedia module ...... no
    SVG module ............. yes
    WebKit module .......... yes
    JavaScriptCore JIT ..... To be decided by JavaScriptCore
    Declarative module ..... no
    Support for S60 ........ no
    Symbian DEF files ...... no
    STL support ............ no
    PCH support ............ no
    MMX/3DNOW/SSE/SSE2/SSE3. no/no/no/no/no
    SSSE3/SSE4.1/SSE4.2..... no/no/no
    AVX..................... yes
    Graphics System ........ raster
    IPv6 support ........... yes
    IPv6 ifname support .... yes
    getaddrinfo support .... yes
    getifaddrs support ..... yes
    Accessibility .......... no
    NIS support ............ no
    CUPS support ........... no
    Iconv support .......... no
    Glib support ........... no
    GStreamer support ...... no
    PulseAudio support ..... no
    Large File support ..... yes
    GIF support ............ yes
    TIFF support ........... yes (qt)
    JPEG support ........... yes (qt)
    PNG support ............ yes (qt)
    MNG support ............ yes (qt)
    zlib support ........... yes
    Session management ..... no
    OpenGL support ......... no
    OpenVG support ......... no
    NAS sound support ...... no
    XShape support ......... no
    XVideo support ......... no
    XSync support .......... no
    Xinerama support ....... no
    Xcursor support ........ no
    Xfixes support ......... runtime
    Xrandr support ......... no
    Xrender support ........ yes
    Xi support ............. no
    MIT-SHM support ........ no
    FontConfig support ..... yes
    XKB Support ............ no
    immodule support ....... yes
    GTK theme support ...... no
    OpenSSL support ........ yes (run-time)
    Alsa support ........... no
    ICD support ............ no
    libICU support ......... no
    Use system proxies ..... no


    WARNING: Using static linking will disable the use of dynamically
    loaded plugins. Make sure to import all needed static plugins,
    or compile needed modules into the library.

    Finding project files. Please wait...
    ...
    ...
    ...

    Qt is now configured for building. Just run 'gmake'.
    Once everything is built, you must run 'gmake install'.
    Qt will be installed into /home/fcj/dev/wkqt

    To reconfigure, run 'gmake confclean' and 'configure'.


    real	1m37.850s
    user	1m13.584s
    sys	0m14.920s


    ### MAKE & ~ INSTALL
    ##
    $ time make -j3

    $ make install

    $ ls ../wkqt/lib/
    libjscore.a    libQtCore.prl  libQtNetwork.a    libQtScript.la
    libjscore.prl  libQtGui.a     libQtNetwork.la   libQtScript.prl
    libQtCore.a    libQtGui.la    libQtNetwork.prl  libQtSql.a
    libQtCore.la   libQtGui.prl   libQtScript.a     libQtSql.la
    libQtSql.prl  libQtTest.a    libQtWebKit.la   libQtXmlPatterns.a
    libQtSvg.a    libQtTest.la   libQtWebKit.prl  libQtXmlPatterns.la
    libQtSvg.la   libQtTest.prl  libQtXml.a       libQtXmlPatterns.prl
    libQtSvg.prl  libQtWebKit.a  libQtXml.la      libQtXml.prl
    pkgconfig/

### Building wkhtmltopdf :

    $ git clone git://github.com/antialize/wkhtmltopdf.git

    # Add a missing version tag :
    $ git tag 0.11.0_rc2 1a24c1717f16fffec7dac90e36bef07e616d5058

    $ git checkout -b live tags/0.11.0_rc2

    $ ../wkqt/bin/qmake

    $ make

    $ ls bin/
    libwkhtmltox.so* (symlinks to:) libwkhtmltox.so.0.10.0
    wkhtmltoimage
    wkhtmltopdf


    $ bin/wkhtmltopdf
    You need to specify atleast one input file, and exactly one output file
    Use - for stdin or stdout

    Name:
      wkhtmltopdf 0.10.0 rc2

    Synopsis:
      wkhtmltopdf [GLOBAL OPTION]... [OBJECT]... <output file>

    Document objects:
      wkhtmltopdf is able to put several objects into the output file, an object is
      either a single webpage, a cover webpage or a table of content.  The objects
      are put into the output document in the order they are specified on the
      command line, options can be specified on a per object basis or in the global
      options area. Options from the Global Options section can only be placed in
      the global options area

      A page objects puts the content of a singe webpage into the output document.

      (page)? <input url/file name> [PAGE OPTION]...
      Options for the page object can be placed in the global options and the page
      options areas. The applicable options can be found in the Page Options and
      Headers And Footer Options sections.

      A cover objects puts the content of a singe webpage into the output document,
      the page does not appear in the table of content, and does not have headers
      and footers.

      cover <input url/file name> [PAGE OPTION]...
      All options that can be specified for a page object can also be specified for
      a cover.

      A table of content object inserts a table of content into the output document.

      toc [TOC OPTION]...
      All options that can be specified for a page object can also be specified for
      a toc, further more the options from the TOC Options section can also be
      applied. The table of content is generated via XSLT which means that it can be
      styled to look however you want it to look. To get an aide of how to do this
      you can dump the default xslt document by supplying the
      --dump-default-toc-xsl, and the outline it works on by supplying
      --dump-outline, see the Outline Options section.

    Description:
      Converts one or more HTML pages into a PDF document, using wkhtmltopdf patched
      qt.

    Global Options:
          --collate                       Collate when printing multiple copies
                                          (default)
          --no-collate                    Do not collate when printing multiple
                                          copies
          --copies <number>               Number of copies to print into the pdf
                                          file (default 1)
      -H, --extended-help                 Display more extensive help, detailing
                                          less common command switches
      -g, --grayscale                     PDF will be generated in grayscale
      -h, --help                          Display help
      -l, --lowquality                    Generates lower quality pdf/ps. Useful to
                                          shrink the result document space
      -O, --orientation <orientation>     Set orientation to Landscape or Portrait
                                          (default Portrait)
      -s, --page-size <Size>              Set paper size to: A4, Letter, etc.
                                          (default A4)
      -q, --quiet                         Be less verbose
          --read-args-from-stdin          Read command line arguments from stdin
          --title <text>                  The title of the generated pdf file (The
                                          title of the first document is used if not
                                          specified)
      -V, --version                       Output version information an exit

    Contact:
      If you experience bugs or want to request new features please visit
      <http://code.google.com/p/wkhtmltopdf/issues/list>, if you have any problems
      or comments please feel free to contact me: see
      <http://www.madalgo.au.dk/~jakobt/#about>
