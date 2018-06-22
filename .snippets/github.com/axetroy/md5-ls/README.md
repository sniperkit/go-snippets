## md5 cli tool

[![Build Status](https://travis-ci.org/axetroy/md5.svg?branch=master)](https://travis-ci.org/axetroy/md5)
![License](https://img.shields.io/badge/license-Apache-green.svg)

List of files/directory MD5 value

## Usage

> Make sure $GOPATH/bin have been add to $PATH

> example: export PATH=$PATH:$GOPATH/bin

```bash
go get -v github.com/axetroy/md5-ls
md5-ls

# output the result like this

63221d5cb2a8851fb0369c986171b493  .all-contributorsrc
55e91fcd2c32356ba54d9204d03cf581  .editorconfig
da8aa4de3a27334fd8be2c2fde195de3  .gitattributes
df10ad3a47f0077b52924b69e2a5b20a  .gitignore
22aaefc1eb93b86248308cd266fa8c66  .gpmrc
af03c276ef40718614dc275ec14c5fa1  .travis.yml
7cbce181f3c83967407f0579a475a53a  LICENSE
af21837b0251ff6e08ddeaf4adbac7a3  README.md
bbe5505e463e8af8ffb2e97a731e4e4c  build
daa537d2c4de1efcd0932078dab3f13f  contributing.md
5679eaa2f73385cc25d5512d4eed2995  main.go
```

Print the cli help information

```bash

md5-ls --help

NAME:
   md5-ls - md5-ls [path]

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.0

DESCRIPTION:
   Display the file/dir md5 value

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Contributing

[Contributing Guid](https://github.com/axetroy/md5-ls/blob/master/CONTRIBUTING.md)

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
| [<img src="https://avatars1.githubusercontent.com/u/9758711?v=3" width="100px;"/><br /><sub>Axetroy</sub>](http://axetroy.github.io)<br />[💻](https://github.com/axetroy/md5-ls/commits?author=axetroy) [🐛](https://github.com/axetroy/md5-ls/issues?q=author%3Aaxetroy) 🎨 |
| :---: |
<!-- ALL-CONTRIBUTORS-LIST:END -->

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Faxetroy%2Fmd5-ls.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Faxetroy%2Fmd5-ls?ref=badge_large)
