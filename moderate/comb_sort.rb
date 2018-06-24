SHRINK = 1.25

File.open(ARGV[0]).each_line do |line|
  t = line.split.map(&:to_i)
  n = 0
  gap = t.length
  until t.each_cons(2).all? { |a, b| (a <=> b) <= 0 }
    n += 1
    gap = (gap / SHRINK).to_i if gap > 1
    (0...t.length - gap).each do |i|
      next if t[i] <= t[i + gap]
      t[i], t[i + gap] = t[i + gap], t[i]
    end
  end
  puts n
end
