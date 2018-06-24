File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(' | ')
  t = s[0].split
  m = s[1].to_i
  while t.length > 1
    n = m % t.length - 1
    t.delete_at((n >= 0) ? n : t.length - 1)
  end
  puts t[0]
end
