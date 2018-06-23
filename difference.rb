require 'json'
require 'pry'
require 'date'
require 'json/add/exception'

tnt_results = File.read('prod_rr.json')
tnt_results = JSON.parse(tnt_results)

different_results = []
at_exit do
  File.open('prod_differing_results.json', 'w') do |f|
    f.write(different_results.to_json)
  end
end

#{
#    "address": {
#        "recipient": "Claudia J Arvidson",
#        "company_name": ""
#        "address": "35616 Ravine Circle",
#        "address_two": "",
#        "city": "Farmington Hills",
#        "state": "MI",
#        "zip": "48335",
#        "country": "US",
#    },
#    "order_id": 658990,
#    "ServiceAvailabilityRequest": {
#        "response": [...],
#        "wiredump": [...],
#    },
#    "RateRequest": {
#        "response": [...],
#        "wiredump": [...],
#    },
#}

class Result
  attr_reader :service, :time_in_transit, :delivery_time
  def initialize(hash)
    @service = hash["Service"]
    @time_in_transit = hash["TimeInTransit"]
    @delivery_time = Date.parse(hash["DeliveryTime"])
  end
end

tnt_results.each do |result|
  begin
    next if !result["rr_errored"].nil? || !result["sar_errored"].nil?

    sar = result["ServiceAvailabilityRequest"]["response"].map{|r| Result.new(r)}
    rr = result["RateRequest"]["response"].map{|r| Result.new(r)}


    if sar.size != rr.size
      different_results << result.merge!(reason: "more service availability results than rate request results")
      next
    end

    broke = false
    sar.each do |sa_result|
      rr_result, rr = rr.partition {|r| r.service == sa_result.service}
      rr_result = rr_result[0]

      if rr_result.nil?
        different_results << result.merge!(reason: "no rate request result")
        broke = true
        break
      end

      if rr_result.time_in_transit != sa_result.time_in_transit
        different_results << result.merge!(reason: "time in transit mismatch")
        broke = true
        break
      end

      if !(rr_result.delivery_time === sa_result.delivery_time)
        different_results << result.merge!(reason: "delivery time mismatch")
        broke = true
      end
    end

    next if broke

    if rr.size > 0
      different_results << result.merge!(reason: "more rate request results than service availability results")
    end
  rescue => e
    different_results << result.merge!(d_errored: e)
  end
end

puts "#{different_results.size} records differed"
