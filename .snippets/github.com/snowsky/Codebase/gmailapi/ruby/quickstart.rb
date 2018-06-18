require 'contextio'

contextio = ContextIO.new('o8h5bh5y', '3Uk6rxMQ8igYO0y2')

account = contextio.accounts.where(email: 'hao.1.wang@gmail.com').first

account.messages.where(limit: 5).each do |message|
  puts message.subject
end
