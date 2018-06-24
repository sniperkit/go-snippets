def choose(n, k)
  (1..k).inject(1) { |a, e| a * (n - e + 1) / e }
end

File.open(ARGV[0]).each_line do |line|
  m = line.to_i
  r = 0
  (0...m / 2).each do |i|
    r += (-1)**i * choose(m, i) * choose(11 * m / 2 - 1 - 10 * i, m - 1)
  end
  puts r
end
