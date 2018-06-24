File.open(ARGV[0]).each_line do |line|
  s = line.chomp
  next if s.empty?
  m = n = -1
  loop do
    m = n
    n = s.index(s[-1], n + 1)
    break if n == s.length - 1
  end
  puts m
end
