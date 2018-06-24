File.open(ARGV[0]).each_line do |line|
  s = line.split(';')
  t = s[0].split(',')
  n = s[1].to_i
  n.step(t.length, n) do |c|
    (1..n / 2).each do |i|
      t[c - i], t[c - n + i - 1] = t[c - n + i - 1], t[c - i]
    end
  end
  puts t.join(',')
end
