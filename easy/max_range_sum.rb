File.open(ARGV[0]).each_line do |line|
  s = line.split(';')
  n = s[0].to_i
  t = s[1].split.map(&:to_i)
  m = c = 0
  (0...n).each { |i| c += t[i] }
  m = c if c > m
  (n...t.length).each do |i|
    c = c - t[i - n] + t[i]
    m = c if c > m
  end
  puts m
end
