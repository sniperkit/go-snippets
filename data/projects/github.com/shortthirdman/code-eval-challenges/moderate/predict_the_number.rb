File.open(ARGV[0]).each_line do |line|
  s = line.to_i
  r = 0
  while s > 0
    s &= s - 1
    r += 1
  end
  puts r % 3
end
