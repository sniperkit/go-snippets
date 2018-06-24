File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  n = s.shift
  s.sort!
  u = s[n / 2]
  puts s.inject(0) { |a, e| a + (e - u).abs }
end
