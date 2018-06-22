barakb/create-branch
====================

A docker file that create image that serve webapp that can create or delete git branches from multiple repositories.

### Downlading:

`docker pull barakb/create-branch-docker`

### Configuring:

create a conf dir named 'conf'  with the following files:

`conf/env.sh` with the following conetnt

```bash
  export ClientID="github client id"
  export ClientSecret="github client secrete"
```

`repos.txt`

 A file with one line for each git repository in the format owner/reponame for example:

`barakb/foobar`


### Running

`./run-container.sh conf`
