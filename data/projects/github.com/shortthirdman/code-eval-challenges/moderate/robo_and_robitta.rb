File.open(ARGV[0]).each_line do |line|
  r = 0
  s = line.scan(/([0-9]+)x([0-9]+) \| ([0-9]+) ([0-9]+)/)[0].map(&:to_i)
  while s[1] > s[3]
    r += s[0]
    s = s[1] - 1, s[0], s[1] - s[3], s[2]
  end
  puts r + s[2]
end
