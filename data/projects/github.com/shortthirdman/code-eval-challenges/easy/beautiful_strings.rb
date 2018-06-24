File.open(ARGV[0]).each_line do |line|
  a = Array.new(26, 0)
  r = 0
  i = 26
  line.chomp.downcase.each_char do |x|
    a[x.ord - 97] += 1 if x =~ /[a-z]/
  end
  a.sort!
  while i > 0 && a[i - 1] > 0
    r += i * a[i - 1]
    i -= 1
  end
  puts r
end
