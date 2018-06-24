File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(',')
  puts s[0].end_with?(s[1]) ? 1 : 0 if s.length > 1
end
