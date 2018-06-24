File.open(ARGV[0]).each_line do |line|
  s = line.split('|')
  t = s[0].split.map(&:to_i)
  n = s[1].to_i
  l = 1
  (1..n).each do
    p = l
    l = 0
    (p...t.length).each do |i|
      next if t[i - 1] <= t[i]
      t[i - 1], t[i] = t[i], t[i - 1]
      l = i > 1 ? i - 1 : 2
      break
    end
    break if l == 0
  end
  puts t.join(' ')
end
