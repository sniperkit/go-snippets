# Design

## Overview
There are a few basic components: servers, hunters, instrumenters, and builders. The servers oversee everything and provide the backend API and frontend website. The hunters are responsible for finding the requested Dockerfile. The instrumenters modify Dockerfiles for the requested architectures and tar everything up. The builders build the images and push to the registry.

Users are allowed to request that a certain repository be built. When this happens, the hunters will attempt to find a Dockerfile for that repository. If one is found, it will configured and sent to the builders. If a relevant Dockerfile is not found, the user will be asked for more information.

By default, built images will reside in the `cobsarch` namespace, where `arch` is replaced with the results of `uname -m`. Repository names will be prepended with their original namespace. For example, if a user requests an ARMv7 build of `dockerfile/rethinkdb`, the resulting repo will be named `cobsarmv7l/dockerfile-rethinkdb`.

## Servers
### Frontend
The basic interface asks for a repo name which may be from either on Docker Hub or GitHub. There is then a set of checkboxes for each of the supported architectures. When the user hits submit, a hunter begins searching the Hubs for the Dockerfile. If one isn't found, an alert is sent back to the user.
Advanced options may include different registries (need to handle auth), different notifications, etc.

### Backend
The backend exposes the API and manages all messaging, logging, etc. It is the only piece that talks directly to the databases.

## Hunters
The hunters try to map the repo name to a Dockerfile and all relevant files. If given a Docker Hub repo, the hunter will first check if it is an Automated Build repo. If so, the Dockerfile is available. If not, the hunter will grab the description and attempt to search for links to the Dockerfile in there. The GitHub hunter searches the repo for the Dockerfile. Once found, the Dockerfiles are then passed to the Instrumenters.

## Instrumenters
The instrumenters parse the Dockerfile and rewrite relevant portions. Any ADDed local files are goroutined off for acquisition. The backend is asked if a mapping exists from the FROM image name to one on each target arch. If a mapping is found, the target arch's name is substituted. If one is not found, the backend will alert the hunters.

Once everything is ready, the Dockerfile and extra baggage is tar'd up and pushed to the builders.

## Builders
The builders push the received tarball to Docker. If Docker finishes the build, the image is pushed to the registry.
