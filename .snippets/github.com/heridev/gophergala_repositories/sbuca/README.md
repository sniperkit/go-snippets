# sbuca

Simple But Useful Certificate Authority

When developing, it's always a pain to generate certificate for SSL/TLS usage. `sbuca` is the simple CA that helps you to generate what you just need in a painless way.

Current features:

1. generate a rsa key (no need to connect to the server)
2. generate a certification request (no need to connect to the server)
3. submit the certification request to the sbuca CA server, and get the signed Certificate
4. get CA's certification 

### Video

<http://youtu.be/6d83GZpt0O8>

### ChallangePost

<http://gopher-gala.challengepost.com/submissions/32201-sbuca>


## Installation

    go get github.com/gophergala/sbuca


## Quick Start

Let's tried the hosted sbuca server

Generate a Rsa Key

    sbuca genkey > test.key

Generate a Certificate Request

    sbuca gencsr --key test.key > test.csr

Submit the Certificate Request to the hosted server and get the Certificate

    sbuca submitcsr --host try.sbuca.com:8600 test.csr > test.crt

In case you want to get the Certificate in another computer, you can add `--format id`, then the output will become the id (serial number) of the Certificate

    sbuca submitcsr --host try.sbuca.com:8600 --format id test.csr 

Then you can get the certificate in another computer (I use ID=2 as example here)

    sbuca getcrt --host try.sbuca.com:8600 2 > test.crt

To get CA's certificate

    sbuca getcacrt --host try.sbuca.com:8600 > ca.crt


## Usage

### Run the CA server

To run a CA server, you can use

    sbuca server

It'll generate ca/ca.srl, ca/ca.key, and ca/ca.crt if needed.
The server listens to `0.0.0.0:8600` by default.


If you want to generate the key & certiricate by your own:

    mkdir ca certs
    echo 01 > ca/ca.srl
    openssl genrsa -out ca/ca.key 2048
    openssl req -x509 -new -key ca/ca.key -out ca/ca.crt
    sbuca server


### Generate a RSA Key

    sbuca genkey > test.key

This command is same as

    openssl genrsa -out test.key 2048


### Generate a Certification Request

    sbuca gencsr --key test.key > test.csr

This command is same as

    openssl req -new -key test.key -out test.csr


### Submit the Certification Request to the CA

By default, it'll output the signed Certificate to STDIN 

    sbuca submitcsr --host localhost:8600 test.csr > test.crt

If you want to get the ID instead, you can add `--format id`

    sbuca submitcsr --host localhost:8600 --format id test.csr

We can use this `id` to get the certificate in another computer

In case you want to use curl to submit the csr, it'll output a JSON by default

    curl localhost:8600/certificates -XPOST --data-urlencode csr@test.csr

If you want to download the Certificate instead of the JSON, you can add `?format=file`

    curl localhost:8600/certificates?format=file -XPOST --data-urlencode csr@test.csr > test.crt


### Get the Certificate by ID

    sbuca getcrt --host localhost:8600 [ID] > test.crt

You can also use curl to get the Certificate (I use ID=2 as example)

    curl localhost:8600/certificates/2?format=file > test.crt


### Get CA's Certificate

    sbuca getcacrt --host localhost:8600 > ca.crt

You can also use curl to get it

    curl localhost:8600/ca/certificate?format=file > ca.crt


### For Docker User

    docker pull waitingkuo/sbuca
    docker run sbuca -p 8600:8600 sbuca server

## TODO

1. password protection
2. admin fuctions: delete, delete all, get all, reset, ...
