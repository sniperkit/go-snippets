File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  m = ''
  s.each { |i| m = i if i.length > m.length }
  puts m
end
