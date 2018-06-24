File.open(ARGV[0]).each_line do |line|
  s = line.split(',').map(&:to_i)
  l = 0
  m = s[0]
  s.each do |i|
    m = i + l if i + l > m
    l = i + l > 0 ? i + l : 0
  end
  puts m
end
