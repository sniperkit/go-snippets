File.open(ARGV[0]).each_line do |line|
  n = line.to_i
  if n < 2
    puts n
  else
    r = p = 1
    r, p, n = r + p, r, n - 1 while n > 2
    puts r
  end
end
