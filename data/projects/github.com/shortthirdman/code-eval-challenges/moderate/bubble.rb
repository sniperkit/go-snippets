File.open(ARGV[0]).each_line do |line|
  s = line.split('|')
  t = s[0].split.map(&:to_i)
  n = s[1].to_i
  (0..t.length - n).each do |i|
    t[i, n + 1] = t[i, n + 1].sort
  end
  puts t.join(' ')
end
