def floors(e, s)
  return 1 if s == 1
  return s if e == 1
  return 0 if e == 0 || s == 0
  return floors(s, s) if e > s
  floors(e - 1, s - 1) + floors(e, s - 1) + 1
end

File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  n = 0
  n += 1 while floors(s[0], n) < s[1]
  puts n
end
