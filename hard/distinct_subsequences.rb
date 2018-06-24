def subs(s, t)
  return 1 if t.length == 0
  while s.length > 0
    return subs(s[1..-1], t[1..-1]) + subs(s[1..-1], t) if s[0] == t[0]
    s = s[1..-1]
  end
  0
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(',')
  puts subs(s[0], s[1])
end
