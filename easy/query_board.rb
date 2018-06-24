m = Array.new(256, Hash.new(0))
File.open(ARGV[0]).each_line do |line|
  s = line.split
  s1 = s[1].to_i
  s2 = s[2].to_i if s.length > 1
  case s[0]
  when 'SetRow' then m[s1] = Hash.new(s2)
  when 'SetCol' then (0..255).each { |i| m[i][s1] = s2 }
  when 'QueryRow'
    r = 0
    (0..255).each { |i| r += m[s1][i] }
    puts r
  when 'QueryCol'
    r = 0
    (0..255).each { |i| r += m[i][s1] }
    puts r
  end
end
