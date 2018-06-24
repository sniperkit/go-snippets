File.open(ARGV[0]).each_line do |line|
  l, d, n, *s = line.split.map(&:to_i)
  count = t = 0
  tx = 6 - d
  i = 6
  while i <= l - 6
    if i > tx - d
      i = tx
      if t == n
        tx = l - 6 + d
      else
        tx = s[t]
        t += 1
      end
    else
      count += 1
    end
    i += d
  end
  puts count
end
