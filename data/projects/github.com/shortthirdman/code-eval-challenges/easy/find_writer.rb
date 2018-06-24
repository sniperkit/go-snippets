File.open(ARGV[0]).each_line do |line|
  s = line.split('|')
  puts s[1].split.map { |x| s[0][x.to_i - 1] }.join if s[1]
end
