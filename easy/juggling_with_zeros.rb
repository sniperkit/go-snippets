File.open(ARGV[0]).each_line do |line|
  r = 0
  s = line.chomp.split.map(&:length)
  (s.length / 2).times do |i|
    r *= 2**s[i * 2 + 1]
    r += 2**s[i * 2 + 1] - 1 if s[i * 2] == 2
  end
  puts r
end
