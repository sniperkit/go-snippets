devsearch-concat
================

Concatenate source files from the DevMine source repositories.

The size of a block on hdfs is at least 64MB. For that reason, if we want to run some large computation
with spark or hadoop's MapReduce we need to concatenate small files into bigger ones that are more suitable for hdfs.
 
devsearch-concat will walk throught the GitHub data that has been made available by DevMine's crawld (https://github.com/DevMine/crawld)
and filter out all files that are not text or too large to be human readable code. It will then create tarballs 
at least 128MB in size with those files.

devsearch-concat assumes a directory structure as follows:

```
REPO_ROOT
└── Language Folder
    └── Github User
        └── Repository
```

The repositories can either be normal directories or tar archives.

All the files' paths in the resulting tar archives are relative to REPO_ROOT. 


Build & Run
-----------

```
> sbt assembly
> java -jar target/scala-2.10/devsearch-concat-assembly-1.0.jar [-j=<numJobs>] <REPO_ROOT> <OUTPUT_FOLDER> 
```
