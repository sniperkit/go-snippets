#!/usr/bin/ruby

require 'json'
require 'pry'
require 'json/add/exception'
require 'rest-client'
require 'thread'

addresses = File.read(ARGV[0])
addresses = JSON.parse(addresses)

at_exit do
  File.open("#{ARGV[1]}.json", 'w') do |f|
    f.write(addresses.to_json)
  end
end

#{
#    "recipient": "Claudia J Arvidson",
#    "company_name": ""
#    "address": "35616 Ravine Circle",
#    "address_two": "",
#    "city": "Farmington Hills",
#    "state": "MI",
#    "zip": "48335",
#    "country": "US",
#}

def request_time_in_transit(addr, skips, lock)
  begin
    req = addr["address"].dup

    skip = false
    lock.synchronize {
      if skips.include?(req)
        skip = true
      else
        skips << req
      end
    }
    return if skip

    req["name"] = req.delete("recipient")
    req["company"] = req.delete("company_name")
    req["address_line_1"] = req.delete("address")
    req["address_line_2"] = req.delete("address_two")
    req["postal_code"] = req.delete("zip")
    req["wiredump"] = true
    req["ship_date"] = "2017-08-29"

    req = URI.encode_www_form(req)
    resp = RestClient.get("localhost:8080/shipping_quote?#{req}")

    addr.merge!({ARGV[1] => JSON.parse(resp)})
  rescue => e
    addr["#{ARGV[1]}_error"] = e
  end
end

q = Queue.new
addresses.each {|a| q.push(a)}
already_run = []
lock = Mutex.new

threads = (1..10).map do
  Thread.new do
    until q.empty?
      addr = q.pop(true) rescue nil
      request_time_in_transit(addr, already_run, lock) if addr
    end
  end
end

threads.each {|t| t.join}

