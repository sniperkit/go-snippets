File.open(ARGV[0]).each_line do |line|
  s = line.split(' | ')
  r = 0
  (1...s.length).each do |x|
    (1...s[0].length).each do |y|
      c = [s[x][y], s[x][y - 1], s[x - 1][y], s[x - 1][y - 1]].sort.join
      r += 1 if c == 'cdeo'
    end
  end
  puts r
end
