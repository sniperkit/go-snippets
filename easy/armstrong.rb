File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split('').map(&:to_i)
  puts s.inject { |a, e| a + e**s.length } == line.to_i ? 'True' : 'False'
end
