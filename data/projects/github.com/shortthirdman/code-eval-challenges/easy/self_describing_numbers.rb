File.open(ARGV[0]).each_line do |line|
  a = []
  n = line.to_i
  line.chomp.each_char do
    a << 0
  end
  m = n
  selfn = true
  while m > 0
    v = m % 10
    if v < a.length
      a[v] += 1
    else
      selfn = false
      break
    end
    m /= 10
  end
  if selfn
    r = 0
    (0...a.length).each do |i|
      r = r * 10 + a[i]
    end
    selfn = false if n != r
  end
  puts selfn ? 1 : 0
end
