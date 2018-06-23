# phUML

phUML is fully automatic UML diagram generator. It takes arbitrary object oriented code written in php5 and creates fully blown class diagrams of it.

-----

[![Docker Stars](https://img.shields.io/docker/stars/dockette/phuml.svg?style=flat)](https://hub.docker.com/r/dockette/phuml/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/phuml.svg?style=flat)](https://hub.docker.com/r/dockette/phuml/)

## Discussion / Help

[![Join the chat](https://img.shields.io/gitter/room/dockette/dockette.svg?style=flat-square)](https://gitter.im/dockette/dockette?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Usage

[![Docker Stars](https://img.shields.io/docker/stars/dockette/phuml.svg?style=flat)](https://hub.docker.com/r/dockette/phuml/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/phuml.svg?style=flat)](https://hub.docker.com/r/dockette/phuml/)

```
docker run dockette/phuml
```

You should see:

```
phUML Version 0.2 (Jakob Westhoff <jakob@php.net>)
Usage: phuml [-h|-l] [OPTIONS] <DIRECTORY> <PROCESSOR> [PROCESSOR OPTIONS] ... <OUTFILE>

Commands:
    -h      Display this help text
    -l      List all available processors

Options:
    -r      Scan given directorie recursively

Example:
    phuml -r ./ -graphviz -createAssociations false -neato out.png

    This example will scan the current directory recursively for php files.
    Send them to the "dot" processor which will process them with the option
    "createAssociations" set to false. After that it will be send to the neato
    processor and saved to the file out.png
```

## Tips

- http://dasunhegoda.com/class-diagram-from-php-code-using-phuml/867/