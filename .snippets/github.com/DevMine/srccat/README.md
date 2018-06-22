# srccat: concatenate repositories into large tar archive

`srccat` reads files from source code repositories, stored either in a folder or
a tar archive, and concatenate them into large tar archives.

This is especially useful for processing on
[Hadoop](http://hadoop.apache.org/) or [Spark](https://spark.apache.org/).

The size of a block on Hadoop Distributed File System (HDFS) is at least 64MB.
Hence, intensive jobs with Spark or Hadoop's MapReduce, perform better if
the small files are concatenated into more suitable larger ones.

`srccat` walks through the files of repositories and filters out all files
that are not text or too large to be human readable code. It then creates
tar archives of a size of at least 128MB with those files.
All the files paths in the resulting tar archives are relative to `REPO_ROOT`.

`srccat` assumes the following directory structure, which is the one used by
[crawld](http://devmine.ch/doc/crawld/):

```
REPO_ROOT
└── Language Folder
    └── Github User
        └── Repository
```

Build & Run
-----------

```
make build
java -jar srccat.jar [-j=<numJobs>] <REPO_ROOT> <OUTPUT_FOLDER>
```
