def contains(s, p)
  until s.empty? || p.empty?
    return false if s[0] > p[0]
    p = p[1..-1] if s[0] == p[0]
    s = s[1..-1]
  end
  p.empty?
end

File.open(ARGV[0]).each_line do |line|
  r = []
  s = line.chomp.split(' | ')
  t = s[0].split
  p = s[1].chars.sort.join
  u = t.map { |i| i.chars.sort.join }
  t.each_with_index { |i, ix| r << i if contains(u[ix], p) }
  puts r.empty? ? 'False' : r.join(' ')
end
