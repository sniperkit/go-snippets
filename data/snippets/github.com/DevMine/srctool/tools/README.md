# Tools

## gvm_cross.sh

This script simply installs cross compilers for multiple plateform (linux,
freebsd, openbsd, netbsd, dragonfly, darwin) and multiple architectures
(32 and 64 bits).

This scripts requires [gvm](https://github.com/moovweb/gvm) to be installed.
However, the authors seem to not care anymore about pull requests and Go 1.4
does not work due to some changes in the source organization. You have to apply
[this patch](https://github.com/tianon/gvm/commit/085833904f596aa4347244c46e1bde6e259bac99).
