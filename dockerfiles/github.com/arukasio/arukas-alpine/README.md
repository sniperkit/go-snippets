arukas-alpine
=============

Getting start in Arukas with Alpine

* Arukas Website: [https://arukas.io/](https://arukas.io/en/)
* Dockerfile: https://github.com/arukasio/arukas-alpine

## Running locally

#### Public key authentication
```
$ git clone git@github.com:arukasio/arukas-alpine.git
$ cd arukas-alpine
$ docker build --no-cache --tag arukasio/alpine .
$ docker run -d -e AUTHORIZED_KEYS="$(cat ~/.ssh/id_rsa.pub)" -P arukasio/alpine
```

#### username/password
If you want to use your original password instead of the default one: "root", you can
set the environment variable ROOT_PASS to your specific password when running the container:
```
$ git clone git@github.com:arukasio/arukas-alpine.git
$ cd arukas-alpine
$ docker build --no-cache --tag arukasio/alpine .
$ docker run -d -e ROOT_PASS="#MyPassW0rd#" -P arukasio/alpine
```
And now you can ssh as root on the container’s IP address  on some port of Docker daemon's host IP address.

## Deploying to Arukas for Arukas User Controll Panel

[Sign-in the Arukas User Controll Panel](https://app.arukas.io/),

#### Public key authentication
```
Image: arukasio/alpine
Instances: 1
Memory: 256MB
Port: 22
ENV: AUTHORIZED_KEYS = ${YOUR_AUTHORIZED_KEYS}
```

#### username/password

```
Image: arukasio/alpine
Instances: 1
Memory: 256MB
Port: 22
ENV: ROOT_PASS = #MyPassW0rd#
```

## Deploying to Arukas for Arukas CLI

[Install the Arukas CLI](https://github.com/arukasio/cli),

or If you have docker installed already,

#### Public key authentication
```
$ docker run \
    --rm \
    -e ARUKAS_JSON_API_TOKEN=${ARUKAS_API_TOKEN} \
    -e ARUKAS_JSON_API_SECRET=${ARUKAS_SECRET_SECRET} \
    arukasio/arukas run --instances=1 --mem=256 --envs AUTHORIZED_KEYS="$(cat ~/.ssh/id_rsa.pub)" --ports=22:tcp arukasio/alpine

```

#### username/password
For demonstration purpose, we’ll create a strong random password in here.
(We strongly recommend your to use a strong password in order to protect your application from crackers and other malicious attacks.)

```
$ ROOT_PASS=$(openssl rand -base64 12)
```

```
$ docker run \
    --rm \
    -e ARUKAS_JSON_API_TOKEN=${ARUKAS_API_TOKEN} \
    -e ARUKAS_JSON_API_SECRET=${ARUKAS_API_SECRET} \
    arukasio/arukas run --instances=1 --mem=256 --ports=22:tcp arukasio/alpine -e ROOT_PASS=${ROOT_PASS}
```

## License

This project is licensed under the terms of the Apache v2 license.
