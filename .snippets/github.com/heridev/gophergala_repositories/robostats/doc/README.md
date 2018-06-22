# Robotstats

This is the internal documentation for Robostats.

## ERD

![erd](https://cloud.githubusercontent.com/assets/385670/5893746/ad786f70-a4b5-11e4-8cf9-5ef50a01de7c.png)

## API endpoints

### POST /user/login (application/x-www-form-urlencoded)

Successful login example:

```sh
curl api.dev.robostats.io/user/login -d "email=user@example.com&password=pass" --verbose
...
> POST /user/login HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-Length: 36
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 36 out of 36 bytes
< HTTP/1.1 200 OK
< Content-Length: 140
< Content-Type: application/json; charset=utf-8
< Date: Sat, 24 Jan 2015 20:06:21 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "access_token": "JKB54JAKcNubIwrcOlukdSQhpZE2Am1ps1tqtlfF",
  "token_type": "bearer",
  "user_id": "54c3fb0960d71e4c5c000007"
}
```

Failed login example:

```sh
curl api.dev.robostats.io/user/login -d "email=user@example.com&password=fail" --verbose
...
< HTTP/1.1 401 Unauthorized
< Content-Length: 13
< Content-Type: text/plain
< Date: Sat, 24 Jan 2015 20:07:57 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
Unauthorized
```

## POST /users (application/json)

Creates an user.

Successful request:

```sh
curl api.dev.robostats.io/user -H "Content-type: application/json" -X POST -d '{"user": {"email": "foo", "password": "pass"}}' --verbose
...
> POST /user HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Content-Length: 46
>
* upload completely sent off: 46 out of 46 bytes
< HTTP/1.1 201 Created
< Content-Length: 3
< Content-Type: text/plain
< Date: Sat, 24 Jan 2015 20:37:07 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "user": {
    "id": "54c4106160d71e5b67000003",
    "email": "fooex",
    "password": "",
    "created_at": "2015-01-24T15:36:32.996882458-06:00",
    "session": {
      "user_id": "54c4106160d71e5b67000003",
      "token": "2yFQe8doL2W3FptLNo5DVMxEDNOuj6NVxFM3HOB2",
      "created_at": "2015-01-24T15:36:33.120703271-06:00"
    }
  }
}
```

Failed request:

```sh
 curl api.dev.robostats.io/user -H "Content-type: application/json" -X POST -d '{"user": {"email": "foo", "password": "pass"}}' --verbose
> POST /user HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Content-Length: 46
>
* upload completely sent off: 46 out of 46 bytes
< HTTP/1.1 422 status code 422
< Content-Length: 48
< Content-Type: application/json
< Date: Sat, 24 Jan 2015 20:40:39 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "errors": [
    "User already exists."
  ]
}
```

## GET /users/:id (application/json)

```
curl api.dev.robostats.io/users/123 -H "Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG" -X GET --verbose
> GET /user/me HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 287
< Content-Type: application/json
< Date: Sat, 24 Jan 2015 23:01:14 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "user": {
    "id": "54c4239260d71e66d5000001",
    "email": "user",
    "password": "",
    "created_at": "2015-01-24T16:58:26.63-06:00",
    "session": {
      "user_id": "54c4239260d71e66d5000001",
      "token": "",
      "created_at": "2015-01-24T16:58:45.372-06:00"
    }
  }
}
```

## POST /device_classes

```
curl api.dev.robostats.io/device_classes -H "Content-type: application/json" -H "Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG" -X POST -d '{"device_classes": {"name": "Class name"}}' --verbose
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 9000 (#0)
> POST /device_classes HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG
> Content-Length: 40
>
* upload completely sent off: 40 out of 40 bytes
< HTTP/1.1 201 Created
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 249
< Content-Type: application/json
< Date: Sat, 24 Jan 2015 23:24:05 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "device_class": {
    "id": "54c4299560d71e6aed000002",
    "user_id": "54c4239260d71e66d5000001",
    "name": "Class name",
    "api_key": "4ZGvOWMSPN2k5M0QDneXLIGRVh2N1m4H9aXPKhKe",
    "created_at": "2015-01-24T17:24:05.128402145-06:00"
  }
}
```

## GET /device_classes

```
curl api.dev.robostats.io/device_classes -H "Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG" -X GET --verbose
> GET /device_classes HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 506
< Content-Type: application/json
< Date: Sat, 24 Jan 2015 23:38:14 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "device_classes": [
    {
      "id": "54c4298960d71e6aed000001",
      "user_id": "54c4239260d71e66d5000001",
      "name": "Class name",
      "api_key": "KfYw0ro4sqPtACpOlj5Xyd7ohTdhU0WmCtinYEC1",
      "created_at": "2015-01-24T17:23:53.398-06:00"
    },
    {
      "id": "54c4299560d71e6aed000002",
      "user_id": "54c4239260d71e66d5000001",
      "name": "Class name",
      "api_key": "4ZGvOWMSPN2k5M0QDneXLIGRVh2N1m4H9aXPKhKe",
      "created_at": "2015-01-24T17:24:05.128-06:00"
    }
  ]
}
```

## GET /device_classes/:id (application/json)

```
curl api.dev.robostats.io/device_classes/54c4299560d71e6aed000002 -H "Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG" -X GET --verbose
> GET /device_classes/54c4299560d71e6aed000002 HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 243
< Content-Type: application/json
< Date: Sat, 24 Jan 2015 23:47:50 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "device_class": {
    "id": "54c4299560d71e6aed000002",
    "user_id": "54c4239260d71e66d5000001",
    "name": "Class name",
    "api_key": "4ZGvOWMSPN2k5M0QDneXLIGRVh2N1m4H9aXPKhKe",
    "created_at": "2015-01-24T17:24:05.128-06:00"
  }
}
```

## DELETE /device_classes/:id

```
curl api.dev.robostats.io/device_classes/54c4299560d71e6aed000002 -H "Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG" -X DELETE --verbose
> DELETE /device_classes/54c4299560d71e6aed000002 HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Authorization: Bearer itdXOKTP9U16B2wtgW1hgpMp0xHKfAkjkCSBKwSG
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 3
< Content-Type: text/plain
< Date: Sat, 24 Jan 2015 23:53:39 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
OK
```

## GET /device_instances

Returns all instances that belong to the user.

Optional GET parameters:

* device_class_id: return only device instances related to this class id.

```
curl api.dev.robostats.io/device_instances -H "Content-type: application/json" -X GET -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071"
> GET /device_instances HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 1314
< Content-Type: application/json
< Date: Sun, 25 Jan 2015 13:35:47 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block

{
  "deviceInstances": [
    {
      "id": "54c4ebfe60d71e7c8400000a",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000008",
      "data": {
        "serial_number": "VQYM1RI5"
      },
      "created_at": "2015-01-25T07:13:34.881-06:00"
    },
    {
      "id": "54c4ebfe60d71e7c8400000b",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000004",
      "data": {
        "serial_number": "BSF4IM9W"
      },
      "created_at": "2015-01-25T07:13:34.882-06:00"
    }
  ]
}
```

## GET /device_instances/:id

Returns the instance that matches the given ID.

```
curl api.dev.robostats.io/device_instances/54c4ebfe60d71e7c8400000e -H "Content-type: application/json" -X GET -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose 2>log.txt
> GET /device_instances/54c4ebfe60d71e7c8400000e HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 259
< Content-Type: application/json
< Date: Sun, 25 Jan 2015 13:37:47 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<

{
  "deviceInstance": {
    "id": "54c4ebfe60d71e7c8400000e",
    "user_id": "54c4ebfe60d71e7c84000001",
    "class_id": "54c4ebfe60d71e7c84000003",
    "data": {
      "serial_number": "HZXTL9LS"
    },
    "created_at": "2015-01-25T07:13:34.883-06:00"
  }
}
```

## DELETE /device_instances/:id

Removes the instance that matches the given ID.

```
curl api.dev.robostats.io/device_instances/54c4ebfe60d71e7c8400000e -H "Content-type: application/json" -X DELETE -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 9000 (#0)
> DELETE /device_instances/54c4ebfe60d71e7c8400000e HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 3
< Content-Type: text/plain
< Date: Sun, 25 Jan 2015 13:38:24 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
OK
```

## GET /device_sessions

Returns all sessions that belong to the user.

```
curl api.dev.robostats.io/device_sessions -H "Content-type: application/json" -X GET -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose
> GET /device_sessions HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Type: application/json
< Date: Sun, 25 Jan 2015 13:42:28 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
< Transfer-Encoding: chunked
<

{
  "deviceSessions": [
    {
      "id": "54c4ebfe60d71e7c8400000f",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000008",
      "instance_id": "54c4ebfe60d71e7c8400000a",
      "session_key": "1thfSqatoPhU2sFQcYMt",
      "data": null,
      "start_time": "2015-01-26T06:09:48.76-06:00",
      "end_time": "2015-02-04T01:36:44.07-06:00",
      "created_at": "2015-01-25T07:13:34.883-06:00"
    },
    {
      "id": "54c4ebfe60d71e7c84000010",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000005",
      "instance_id": "54c4ebfe60d71e7c8400000d",
      "session_key": "h3d7MhuY1SkELdSJwiWo",
      "data": null,
      "start_time": "2015-01-05T23:21:00.806-06:00",
      "end_time": "2015-01-17T21:46:05.884-06:00",
      "created_at": "2015-01-25T07:13:34.884-06:00"
    }
	]
}
```

## GET /device_sessions/:id

Returns the session that matches the given ID.

```
curl api.dev.robostats.io/device_sessions/54c4ebfe60d71e7c8400000f -H "Content-type: application/json" -X GET -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 9000 (#0)
> GET /device_sessions/54c4ebfe60d71e7c8400000f HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 409
< Content-Type: application/json
< Date: Sun, 25 Jan 2015 13:45:08 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
{
  "deviceSession": {
    "id": "54c4ebfe60d71e7c8400000f",
    "user_id": "54c4ebfe60d71e7c84000001",
    "class_id": "54c4ebfe60d71e7c84000008",
    "instance_id": "54c4ebfe60d71e7c8400000a",
    "session_key": "1thfSqatoPhU2sFQcYMt",
    "data": null,
    "start_time": "2015-01-26T06:09:48.76-06:00",
    "end_time": "2015-02-04T01:36:44.07-06:00",
    "created_at": "2015-01-25T07:13:34.883-06:00"
  }
}
```

## DELETE /device_instances/:id

```
curl api.dev.robostats.io/device_sessions/54c4ebfe60d71e7c8400000f -H "Content-type: application/json" -X DELETE -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 9000 (#0)
> DELETE /device_sessions/54c4ebfe60d71e7c8400000f HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Length: 3
< Content-Type: text/plain
< Date: Sun, 25 Jan 2015 13:46:21 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
<
OK
```

## GET /device_events

Returns all logs that belong to the user.

```
curl api.dev.robostats.io/device_events -H "Content-type: application/json" -X GET -H "Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071" --verbose

> GET /device_events HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Content-type: application/json
> Authorization: Bearer lapfkPYXWJkhSasV26jD8VN3unMkVF2LvRht2071
>
  0     0    0     0    0     0      0      0 --:--:--  0:00:02 --:--:--     0< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Content-Type: application/json
< Date: Sun, 25 Jan 2015 13:50:48 GMT
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block
< Transfer-Encoding: chunked
<

{
  "deviceLogs": [
    {
      "id": "54c4ebfe60d71e7c84000019",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000003",
      "instance_id": "54c4ebfe60d71e7c8400000c",
      "session_id": "54c4ebfe60d71e7c84000012",
      "data": null,
      "local_time": 0,
      "latlng": [
        18.75536,
        -98.67731
      ],
      "created_at": "2015-01-25T07:13:34.887-06:00"
    },
    {
      "id": "54c4ebfe60d71e7c8400001a",
      "user_id": "54c4ebfe60d71e7c84000001",
      "class_id": "54c4ebfe60d71e7c84000005",
      "instance_id": "54c4ebfe60d71e7c8400000d",
      "session_id": "54c4ebfe60d71e7c84000011",
      "data": null,
      "local_time": 1,
      "latlng": [
        19.97565,
        -99.20467000000001
      ],
      "created_at": "2015-01-25T07:13:34.887-06:00"
    }
	]
}
```

## GET /device_sessions/time_series

Required parameters:

* session_id
* key[]

```
curl 'api.dev.robostats.io/device_sessions/time_series?session_id=54c542d612fa74250100000f&key\[\]=cpu' -X GET -H "Authorization: Bearer W7wDOMQ4ytC0hfI76r3HdBOGSoiZ1lB01vNOLx8T"  --verbose
* Hostname was NOT found in DNS cache
*   Trying 104.236.93.29...
* Connected to api.dev.robostats.io (104.236.93.29) port 80 (#0)
> GET /device_sessions/time_series?session_id=54c542d612fa74250100000f&key[]=cpu HTTP/1.1
> User-Agent: curl/7.37.0
> Host: api.dev.robostats.io
> Accept: */*
> Authorization: Bearer W7wDOMQ4ytC0hfI76r3HdBOGSoiZ1lB01vNOLx8T
>
< HTTP/1.1 200 OK
* Server nginx/1.4.6 (Ubuntu) is not blacklisted
< Server: nginx/1.4.6 (Ubuntu)
< Date: Sun, 25 Jan 2015 19:28:53 GMT
< Content-Type: application/json
< Transfer-Encoding: chunked
< Connection: keep-alive
< Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
< Access-Control-Allow-Origin: *
< Set-Cookie: REVEL_FLASH=; Path=/
< X-Content-Type-Options: nosniff
< X-Frame-Options: SAMEORIGIN
< X-Xss-Protection: 1; mode=block

{
  "time_serie":{
    "steps":[
      "4",
      "9",
      "12",
      "16",
      "18",
      "23",
      "24",
      "28",
      "29",
      "30",
      "33",
      "38",
      "39",
      "40",
      "45",
      "50",
      "53",
      "56",
      "60",
      "63",
      "68",
      "73",
      "74",
      "79",
      "84",
      "86",
      "87",
      "91",
      "94",
      "98",
      "102",
      "107",
      "110",
      "111",
      "116",
      "121",
      "125",
      "128",
      "129",
      "130",
      "132",
      "135",
      "137",
      "139",
      "144",
      "146",
      "149",
      "151",
      "155",
      "160",
      "165",
      "168",
      "173",
      "177",
      "182",
      "184",
      "185",
      "186",
      "191",
      "195",
      "200",
      "202",
      "205",
      "206",
      "210",
      "212",
      "217",
      "220",
      "222",
      "226",
      "227",
      "228",
      "231",
      "233",
      "236",
      "237",
      "242",
      "247",
      "249",
      "253",
      "258",
      "262",
      "267",
      "271",
      "272",
      "277",
      "280",
      "283",
      "288",
      "291",
      "292",
      "293",
      "298",
      "300",
      "301"
    ],
    "values":{
      "cpu":[
        0.06580421328544617,
        0.8488441705703735,
        0.10052532702684402,
        0.2614118158817291,
        0.6557889580726624,
        0.6584054827690125,
        0.9775729179382324,
        0.7425039410591125,
        0.43388327956199646,
        0.26217707991600037,
        0.48331668972969055,
        0.2840568721294403,
        0.9878626465797424,
        0.38125765323638916,
        0.640898585319519,
        0.6719640493392944,
        0.6014140248298645,
        0.0981513187289238,
        0.22008708119392395,
        0.5315225720405579,
        0.43042775988578796,
        0.08740442991256714,
        0.004428897984325886,
        0.8683210611343384,
        0.7270494699478149,
        0.3117087781429291,
        0.43723079562187195,
        0.7535176277160645,
        0.23146456480026245,
        0.7297129034996033,
        0.04842108115553856,
        0.3271646201610565,
        0.44170287251472473,
        0.20092080533504486,
        0.9143272638320923,
        0.18377190828323364,
        0.20140109956264496,
        0.658624529838562,
        0.7657321691513062,
        0.005155558697879314,
        0.772403359413147,
        0.6709102392196655,
        0.30460798740386963,
        0.50579434633255,
        0.48936089873313904,
        0.14246414601802826,
        0.2330542653799057,
        0.9285175800323486,
        0.7499865889549255,
        0.1671348214149475,
        0.7868396043777466,
        0.7580482959747314,
        0.13535353541374207,
        0.7484891414642334,
        0.961996853351593,
        0.5624016523361206,
        0.29574349522590637,
        0.12745633721351624,
        0.9174259305000305,
        0.20055775344371796,
        0.07886053621768951,
        0.8275659084320068,
        0.20144103467464447,
        0.4306133985519409,
        0.8462740778923035,
        0.8102183938026428,
        0.7332144379615784,
        0.316532701253891,
        0.10482120513916016,
        0.6120192408561707,
        0.06303984671831131,
        0.26034507155418396,
        0.04752376675605774,
        0.42816299200057983,
        0.9987062215805054,
        0.29561689496040344,
        0.5344931483268738,
        0.2778962254524231,
        0.2912706732749939,
        0.31632083654403687,
        0.4529159367084503,
        0.7359374165534973,
        0.4961923062801361,
        0.3104954957962036,
        0.9962745308876038,
        0.6005687713623047,
        0.9811151623725891,
        0.0733659490942955,
        0.1772964596748352,
        0.1126754954457283,
        0.1941128522157669,
        0.33268168568611145,
        0.11313501745462418,
        0.984341561794281,
        0.026709144935011864
      ]
    },
    "label":{
      "cpu":"cpu"
    }
  }
}
```

