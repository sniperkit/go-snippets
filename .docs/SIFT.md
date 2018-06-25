# Sift 

## Example

```bash
sift --no-zip \
     --no-conf \
     --output-unixpath \
     --git \
     --binary-skip \
     --no-color \
     --group \
     --dirs=vendor,Godeps,node_modules,bower_components \
	 --file-matches==vendor,Godeps,node_modules,bower_components \
     --targets
```
# --output=clean_dirs.txt

## Documentation
```bash
Usage:
		sift   [OPTIONS]  PATTERN  [FILE|PATH|tcp://HOST:PORT]...
		sift [OPTIONS] [-e PATTERN | -f FILE]  [FILE|PATH|tcp://HOST:PORT]...
		sift [OPTIONS] --targets [FILE|PATH]...

OPTIONS
       --binary-skip
              skip files that seem to be binary

       -a, --binary-text
              process files that seem to be binary as text

       --blocksize=
              blocksize in bytes (with optional suffix K|M)

       --color
              enable colored output (default: auto)

       --no-color
              disable colored output

       -C, --context=NUM
              show NUM context lines

       -A, --context-after=NUM
              show NUM context lines after match

       -B, --context-before=NUM
              show NUM context lines before match

       -j, --cores=
              limit used CPU Cores (default: 0 = all)

       -c, --count
              print count of matches per file

       --dirs=GLOB
              recurse only into directories whose name matches GLOB

       --err-show-line-length
              show all line length errors

       --err-skip-line-length
              skip line length errors

       --exclude-dirs=GLOB
              do not recurse into directories whose name matches GLOB

       -x, --ext=
              limit search to specific file extensions (comma-separated)

       -X, --exclude-ext=
              exclude specific file extensions (comma-separated)

       --files=GLOB
              search only files whose name matches GLOB

       --exclude-files=GLOB
              do not select files whose name matches GLOB while recursing

       --path=PATTERN
              search only files whose path matches PATTERN

       --ipath=PATTERN
              search only files whose path matches PATTERN (case insensitive)

       --exclude-path=PATTERN
              do not search files whose path matches PATTERN

       --exclude-ipath=PATTERN
              do  not  search  files whose path matches PATTERN (case insensi-
              tive)

       -t, --type=
              limit  search  to  specific  file  types  (comma-separated,  see
              --list-types)

       -T, --no-type=
              exclude specific file types (comma-separated, --list-types)

       -l, --files-with-matches
              list files containing matches

       -L, --files-without-match
              list files containing no match

       --follow
              follow symlinks

       --git  
              respect .gitignore files and skip .git directories

       --group
              group output by file (default: off)

       --no-group
              do not group output by file

       -i, --ignore-case
              case insensitive (default: off)

       -I, --no-ignore-case
              disable case insensitive

       -s, --smart-case
              case  insensitive  unless  pattern contains uppercase characters
              (default: off)

       -S, --no-smart-case
              disable smart case

       --no-conf
              do not load config files

       -v, --invert-match
              select non-matching lines

       --limit=NUM
              only show first NUM matches per file

       -Q, --literal
              treat pattern as literal, quote meta characters

       -m, --multiline
              multiline parsing (default: off)

       -M, --no-multiline
              disable multiline parsing

       --only-matching
              only show the matching part of a line

       -o, --output=FILE|tcp://HOST:PORT
              write output to the specified file or network connection

       --output-limit=
              limit output length per found match

       --output-sep=
              output separator (default: "\n")

       --output-unixpath
              output file paths in unix format ('/' as path separator)

       -e, --regexp=PATTERN
              add pattern PATTERN to the search

       -f, --regexp-file=FILE
              search for patterns contained in FILE (one per line)

       --print-config
              print config for loaded configs + given command line arguments

       -q, --quiet
              suppress output, exit with return code  zero  if  any  match  is
              found

       -r, --recursive
              recurse into directories (default: on)

       -R, --no-recursive
              do not recurse into directories

       --replace=
              replace  numbered or named (?P<name>pattern) capture groups. Use
              ${1}, ${2}, $name, ... for captured submatches

       --filename
              enforce printing the filename before results (default: auto)

       --no-filename
              disable printing the filename before results

       -n, --line-number
              show line numbers (default: off)

       -N, --no-line-number
              do not show line numbers

       --column
              show column numbers

       --no-column
              do not show column numbers

       --stats
              show statistics

       --targets
              only list selected files, do not search

       --list-types
              list available file types

       -V, --version
              show version and license information

       -w, --word-regexp
              only match on ASCII word boundaries

       --write-config
              save config for loaded configs + given command line arguments

       -z, --zip
              search content of compressed .gz files (default: off)

       -Z, --no-zip
              do not search content of compressed .gz files

   File Condition options:
       --file-matches=PATTERN
              only show matches if file also matches PATTERN

       --line-matches=NUM:PATTERN
              only show matches if line NUM matches PATTERN

       --range-matches=X:Y:PATTERN
              only show matches if lines X-Y match PATTERN

       --not-file-matches=PATTERN
              only show matches if file does not match PATTERN

       --not-line-matches=NUM:PATTERN
              only show matches if line NUM does not match PATTERN

       --not-range-matches=X:Y:PATTERN
              only show matches if lines X-Y do not match PATTERN

   Match Condition options:
       --preceded-by=PATTERN
              only show matches preceded by PATTERN

       --followed-by=PATTERN
              only show matches followed by PATTERN

       --surrounded-by=PATTERN
              only show matches surrounded by PATTERN

       --preceded-within=NUM:PATTERN
              only show matches preceded by PATTERN within NUM lines

       --followed-within=NUM:PATTERN
              only show matches followed by PATTERN within NUM lines

       --surrounded-within=NUM:PATTERN
              only show matches surrounded by PATTERN within NUM lines

       --not-preceded-by=PATTERN
              only show matches not preceded by PATTERN

       --not-followed-by=PATTERN
              only show matches not followed by PATTERN

       --not-surrounded-by=PATTERN
              only show matches not surrounded by PATTERN

       --not-preceded-within=NUM:PATTERN
              only show matches not preceded by PATTERN within NUM lines

       --not-followed-within=NUM:PATTERN
              only show matches not followed by PATTERN within NUM lines

       --not-surrounded-within=NUM:PATTERN
              only show matches not surrounded by PATTERN within NUM lines

   Help Options:
       -h, --help
              Show this help message
```