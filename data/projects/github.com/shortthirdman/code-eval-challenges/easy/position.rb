File.open(ARGV[0]).each_line do |line|
  s = line.split(',').map(&:to_i)
  puts s[0] & 1 << s[1] - 1 > 0 == s[0] & 1 << s[2] - 1 > 0
end
