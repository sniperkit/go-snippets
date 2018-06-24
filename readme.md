# git-ledger [![Build Status](https://travis-ci.org/git-hook/git-ledger.svg?branch=master)](https://travis-ci.org/git-hook/git-ledger)

> Track your local git clones

**git-ledger** acts as a simple key-value store, remembering the
location of git-repositories on local filesystems.

## Install

To install the library

```bash
go get github.com/git-hook/git-ledger
```

To install the binary

```bash
go get github.com/git-hook/git-ledger/cmd/git-ledger
git ledger add .
```

## API

### add [path]

Start tracking an existing repository.

### find [slug]

Print the location of a tracked repository.

### ls

Print all tracked repositories.

### rm [path-or-slug]

Stop tracking an existing repository.

## Related

- [grit](https://github.com/jmalloc/grit)

## License

MIT © Eric Crosson
