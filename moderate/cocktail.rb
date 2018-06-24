File.open(ARGV[0]).each_line do |line|
  s = line.split('|')
  t = s[0].split.map(&:to_i)
  n = [s[1].to_i, t.length / 2].min
  (1..n).each do |k|
    (k..t.length - k).each do |i|
      t[i - 1], t[i] = t[i], t[i - 1] if t[i - 1] > t[i]
    end
    (t.length - k - 1).downto(k) do |i|
      t[i - 1], t[i] = t[i], t[i - 1] if t[i - 1] > t[i]
    end
  end
  puts t.join(' ')
end
