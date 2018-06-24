stats = [[2, 4, 8, 6], [3, 9, 7, 1], [4, 6], [5], [6], [7, 9, 3, 1],
         [8, 4, 2, 6], [9, 1]]

File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  s[0] -= 2
  m = s[1] / stats[s[0]].length
  r = Array.new(10, 0)
  stats[s[0]].each { |x| r[x] += m }
  (0...(s[1] % stats[s[0]].length)).each { |x| r[stats[s[0]][x]] += 1 }
  9.downto(0) { |x| r[x] = x.to_s + ': ' + r[x].to_s }
  puts r.join(', ')
end
