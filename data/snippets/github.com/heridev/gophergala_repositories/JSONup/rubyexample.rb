#!/usr/bin/ruby

require 'rest-client'

url = "http://jsonup.com/push/yrd3jy30udi"



while(true) do
  r = rand(4)
  amt = rand(99)
  json = [
    {
      name: "ruby.test#{r}",
      status: 'UP',
      value: "#{amt}"
    },
    {
      name: "ruby.b#{r}",
      status: 'DOWN',
      value: "#{amt/2}"
    },
  ].to_json

  RestClient.post(url, json)

end
