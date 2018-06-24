def pos(l, r)
  (l + 'a'.ord).chr + (r + '1'.ord).chr
end

File.open(ARGV[0]).each_line do |line|
  l = line[0].ord - 'a'.ord
  r = line[1].ord - '1'.ord
  m = []
  m << pos(l - 2, r - 1) if l > 1 && r > 0
  m << pos(l - 2, r + 1) if l > 1 && r < 7
  m << pos(l - 1, r - 2) if l > 0 && r > 1
  m << pos(l - 1, r + 2) if l > 0 && r < 6
  m << pos(l + 1, r - 2) if l < 7 && r > 1
  m << pos(l + 1, r + 2) if l < 7 && r < 6
  m << pos(l + 2, r - 1) if l < 6 && r > 0
  m << pos(l + 2, r + 1) if l < 6 && r < 7
  puts m.join(' ')
end
