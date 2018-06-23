# CoBS: The Container Build Service
CoBS builds container images for alternative hardware architectures, such as ARM and POWER. It can also be used as a regular build service, similar to [Docker Hub's automated builds](http://docs.docker.com/docker-hub/builds/) and [CoreOS's Quay](https://quay.io).

# Inspiration
Initial inspiration came from a question in #google-containers asking about Kubernetes support on ARM. I had recently acquired an ODROID-C1, so I decided to try it out. And that's where things went wrong.

Even though Docker images have an architecture property attached to them, this is currently ignored when doing a `docker run`. Since most images on Docker Hub are built for x86_64, trying to run them on ARM will result in an exec format error.

CoBS aims to bridge the gap between x86_64 images on Docker Hub and alternative hardware architectures.

# Target Audience
The main target audience are users who want to run containers on alternative architectures, such as the Raspberry Pi. A secondary audience includes users who want to perform their x86_64 container image builds on premises.

# Current Status
Backend API and the image builders are working, but there currently is no frontend.

Build status is currently not reported. You can manually check if your requested image was built by checking the [cobsarmv7l organization on Docker Hub](https://registry.hub.docker.com/repos/cobsarmv7l/).

Builders for x86_64 and armv7l are currently active. An armv6 builder (Raspberry Pi) will hopefully be added soon.

## Try it out
`$ curl -X POST -F repository=redis -F arch=armv7l http://cobs.aas.io/search`
`http://cobs.aas.io/api/v1/info/4297016f-393e-45f5-881f-cfd636155988`

`$ curl http://cobs.aas.io/api/v1/info/4297016f-393e-45f5-881f-cfd636155988`
`{"original_name":"redis","new_name":"cobsarmv7l/redis","architecture":"armv7l","tag":"latest"}`

On your ARMv7 machine (such as the ODROID-C1) you can then do a `docker pull cobsarmv7l/redis` to get a version of Redis that will run on your architecture.

# Frequently Asked Questions

## Why would you want to run Docker on a Raspberry Pi?
One example: local database + web server for your IoT sensors, allowing logging even when Internet connectivity is lost.

## How does this work?
You submit a request with a repository name, desired tag, and target architecture. CoBS then searches Docker Hub's Registry for that repository. If a Dockerfile is attached to the repo, CoBS will analyze it and modify the base (FROM) image to an equivalent one that supports the target architecture. If this base image doesn't exist, CoBS will attempt to build it. Once all dependencies are met, a build machine will attempt to build your image. If successful, the image is pushed to Docker Hub.

## What base images are supported?
Currently debian:{wheezy,jessie,sid} and ubuntu{saucy,trusty,utopic,vivid}.

## The repository I want to use doesn't include a Dockerfile.
Unfortunately we can't build a container without one.

## My image needs more files than just the Dockerfile.
Unfortunately I had to drop support for this for the Gala. But it will come soon.

## My Dockerfile is in a repo on GitHub
Initial support exists, but currently isn't active.

## Where's the POWER?
Hasn't been tested yet due to time constraints and the current questionable support of Docker/Go on POWER.

## Does this work under QEMU?
Probably. For now I am using physical hardware.

## Do you support Rocket?
Not yet. 48 hours isn't enough for a solo project :)

## You have some major security holes.
Yup.

## None of my requests are being built.
About that.... Hopefully it's working again by the time you read this.

## There are hard-coded values everywhere!
Umm....time was limited? Don't worry, it will all be cleaned up soon.

## That's not an excuse. And your code is horrible.
Cut me some slack, this is my first time writing Go :)
