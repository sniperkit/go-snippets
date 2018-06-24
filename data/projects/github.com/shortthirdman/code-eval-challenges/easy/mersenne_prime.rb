require 'prime'
File.open(ARGV[0]).each_line do |line|
  puts Prime.take_while { |i| 2**i - 1 < line.to_i } \
    .map { |i| 2**i - 1 }.join(', ')
end
