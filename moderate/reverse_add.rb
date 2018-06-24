File.open(ARGV[0]).each_line do |line|
  a = line.to_i
  i = 0
  while i < 100
    r = a.to_s.reverse.to_i
    break if a == r
    a += r
    i += 1
  end
  puts a == r ? i.to_s + ' ' + a.to_s : 'not found'
end
