package responses

type Formatter func(status int, body interface{}) error
