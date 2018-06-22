tar fixtures notes
==================

### `tar_sansBase.tgz`

- gzipped.
- produced by gnu tar.
- two entries: `ab` and `cd/` -- a file and dir, respectively
  - note the lack of `./` prefix
- ownership is 7000/7000.  user and group *name* are also included: both are "hash".
- dates are 2015-05-30 14:11:26 -0500

### `tar_withBase.tgz`

- gzipped.
- produced by gnu tar.
- *three* entries: `./`, `./ab`, and `./bc/` -- dir, file, dir.
  - careful: some graphical tar tools may not show you that root dir.
  - note the presense of `./` prefix
- ownership is 7000/7000.  user and group *name* are also included: both are "hash".
- dates are 2015-05-30 14:53:35 -0500
- so, compared to `tar_sansBase.tgz`:
  - they diverge in timestamps, and filenames, and...
  - and the `./` entry: rio normalization will add a base dir placeholder, but it won't have those ownership bits!
  - the `./` prefix should not matter; rio normalization vanishes it completely.

### `tar_kitchenSink.tgz`

- gzipped.
- produced by rio.
- many entires:
  - `./`
  - `./dir/`
  - `./dir/f1`
  - `./deep/`
  - `./deep/tree/`
  - `./deep/tree/f3`
  - `./lnkdangle`
  - `./empty/`
  - `./f2`
  - `./lnkfile`
  - `./lnkdir`
- ownership is mostly 7000:7000, but one file (f2) is 4000:5000.  no usernames.
- dates are various in 2017-09-27.
- a variety of symlinks are included.
