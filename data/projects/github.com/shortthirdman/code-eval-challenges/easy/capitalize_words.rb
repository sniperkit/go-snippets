File.open(ARGV[0]).each_line do |line|
  s = line.split.each { |x| x[0] = x[0].capitalize }
  puts s.join(' ')
end
