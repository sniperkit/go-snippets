File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  r = 0
  (1 << s[0]).upto(s[1]) do |i|
    n = s[0]
    a = i
    while a > 0 && n >= 0
      n -= 1 if a.even?
      a /= 2
    end
    r += 1 if n == 0
  end
  puts r
end
