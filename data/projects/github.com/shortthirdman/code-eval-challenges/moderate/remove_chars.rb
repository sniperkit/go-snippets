File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(', ')
  puts s[0].delete(s[1])
end
