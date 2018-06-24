require 'prime'
File.open(ARGV[0]).each_line do |line|
  puts Prime.take_while { |i| i < line.to_i }.join(',')
end
