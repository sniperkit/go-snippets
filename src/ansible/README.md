# Ansible Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/ansible.svg)
![](https://img.shields.io/docker/build/qmu1/ansible.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/ansible.svg)](https://microbadger.com/images/qmu1/ansible)

An alipine-based dockerfile which runs [Ansible](https://www.ansible.com).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/ansible/Dockerfile)


## Usage

```
docker run --rm -it \
   -v ~/.ssh:/home/docker/.ssh \
   -v $(pwd):/home/docker/ansible \
   qmu1/ansible /bin/sh -c \
   'ansible-playbook --version && ls -la'
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/ansible/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
