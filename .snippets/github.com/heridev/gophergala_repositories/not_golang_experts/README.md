# Not golang experts: Gopherstalker

## API Endpoints

### Registrations

#### POST `/registrations`

Request body:

``` json
{
  "user" : {
    "email" : "email@example.com",
    "password" : "supersecret",
    "password_confirmation" : "supersecret"
  }
}
```

##### Response codes

Status `201`

``` json
{
  "token" : "authtoken"
}
```

Status `422`

Error messages:

* Email has already been taken
* Passwords don't match

```json
{
  "error" : "error message"
}
```

### Sessions

#### POST `/sessions`

Request body:

``` json
{
  "user" : {
    "email" : "email@example.com",
    "password" : "supersecret"
  }
}
```

##### Response codes

Status `201`

``` json
{
  "token" : "authtoken"
}
```

Status `422`

Error messages:

* Invalid email or password

```json
{
  "error" : "error message"
}
```

#### DELETE `/sessions?token=yourauthtoken`

**This request does not require body**

##### Response codes

Status `201`

``` json
{
  "token" : "authtoken"
}
```

Status `404`

Error messages:

* Not found

```json
{
  "error" : "error message"
}
```

### Subscriptions

#### POST `/subscriptions?token=yourauthtoken`

Request body:

``` json
{
  "url" : "http://gophercon.com/"
}
```

##### Response codes

Status `201`

``` json
{
  "Subscription": {
    "Id": 1,
    "UserId": 1,
    "PageId": 1,
    "CreatedAt": "2015-01-25T11:39:50.184738-06:00",
    "UpdatedAt": "2015-01-25T11:39:50.184738-06:00"
  }
}
```

Status `401`

Error messages:

* Invalid session token

```json
{
  "error" : "error message"
}
```

#### GET `/subscriptions?token=yourauthtoken`

** This request does not require body **

##### Response codes

Status `200`

``` json
[
   {
    "Id": 1,
    "UserId": 1,
    "PageId": 1,
    "CreatedAt": "2015-01-25T11:39:50.184738-06:00",
    "UpdatedAt": "2015-01-25T11:39:50.184738-06:00"
  },
  {
    "Id": 2,
    "UserId": 1,
    "PageId": 2,
    "CreatedAt": "2015-01-25T11:39:50.184738-06:00",
    "UpdatedAt": "2015-01-25T11:39:50.184738-06:00"
  }
]
```

Status `401`

Error messages:

* Invalid session token

```json
{
  "error" : "error message"
}
```

#### DELETE `/subscriptions/1?token=yourauthtoken`

** This request does not require body **

##### Response codes

Status `200`

``` json
{
  "message" : "Successfully unsubscribed from subscription: 1"
}
```

Status `401`

Error messages:

* Invalid subscription id
* Invalid session token

```json
{
  "error" : "error message"
}
```
