File.open(ARGV[0]).each_line do |line|
  s = line.split
  l = s[-1].to_i
  puts s[s.length - l - 1] if l < s.length
end
