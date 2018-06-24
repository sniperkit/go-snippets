def pol(o, a, b)
  case o
  when '*'
    return a * b
  when '/'
    return a / b
  end
  a + b
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  n = s.length / 2
  res = s[n].to_f
  (1..n).each do |x|
    res = pol(s[n - x], res, s[n + x].to_f)
  end
  res += 0.001
  puts res.to_i
end
