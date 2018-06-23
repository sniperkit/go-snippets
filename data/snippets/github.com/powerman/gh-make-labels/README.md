# Make labels for GitHub repo

This command-line tool will update given repo's labels to match predefined
labels list.

Check [this repo labels](https://github.com/powerman/gh-make-labels/labels)
as example which labels you'll have after running this tool.

## Installation

```sh
go get github.com/powerman/gh-make-labels
```

## Dependencies

You should have installed and configured
[hub](https://github.com/github/hub) tool (gh-make-labels will use hub's
token to access GitHub API).
