# Scpm

Copy files over ssh protocol to multiple servers, scp multiplexer.

## Demo
```
scpm --in="~/Video/dog-team09012015.mp4" --path="example1.com:/tmp/a" --path="example2.com:/tmp/b"
Start copy /home/gron/Video/dog-team09012015.mp4
example1.com:22 /tmp/a 5.59 MB / 109.99 MB [=>-----------------------------] 5.09 % 1.12 MB/s 1m33s
example2.com:22 /tmp/b 5.69 MB / 109.99 MB [=>-----------------------------] 5.17 % 1.14 MB/s 1m31s
```

## Install
```bash
go get github.com/gophergala/scpm/cmd/scpm
```

## Usage
```
scpm --in="/path/to" \
    --path="user@server0.com:/path/to" \
    --path="user@server1.com:/path/to" \
    --path="user@server2.com:/path/to" \
    --path="user@server3.com:/path/to" \
    --path="user@server4.com:/path/to"

```

## Features
1. Copy once file to once server
2. Copy once file to multiple servers
3. Recursive copy folder to multiple server
4. Multi progress bar
4. Support custom identity key, system identity, or plain password

## TODO
1. Flag config.json