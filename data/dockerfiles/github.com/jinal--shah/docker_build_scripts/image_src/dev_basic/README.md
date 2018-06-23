[1]: https://github.com/jinal--shah/demo-coreos-vagrant-setup "using coreos with vagrant demo"

# Container: dev_basic

## devbox

References to `coreos` and `devbox` are only relevant if you're
using [demo-coreos-vagrant-setup][1]. Otherwise you can safely ignore them.

## BUILD | RUN

* build from Dockerfile in this dir 

  e.g. with tag 0.0.1

  ```
    docker build --no-cache=true --rm --tag dev_basic:0.0.1 .
  ```

* drop in to a bash workspace from this tagged image with `devbox` helper function:

  ```
    # example mount vols /experiments and /my_app using dirs on coreos host
    PROJ_DIR=$HOME/core-01-projects/projects
    devbox $PROJ_DIR/experiments:$PROJ_DIR/my_app dev_basic:0.0.1
  ```

* ... do stuff under /experiments and /my_app then exit shell

  This stops the container.

  Your changes will have persisted under $PROJ_DIR on the coreos host.
  _Also if you used the devbox function, you'll see the commands you ran_
  _in the container are in the `core` user's history on the coreos host._

* restart the same container

  First, get its container id with `docker ps -a`

  Drop back in to the workspace with:

  ```
  docker start -a -i <container_id> 
  ```

  Notice that the same mounted vols are available.

* open 1 or more sessions by running this from different sessions on the coreos host

  ```
  docker exec -it <container_id> /bin/bash
  ```

## DEVBOX FUNCTION

Run `devbox_help` for info, or look at source in .bashrc on coreos host.
