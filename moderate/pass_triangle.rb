a = []
File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split.map(&:to_i)
  s[0] += a[0] if a.length > 0
  s[-1] += a[-1] if a.length > 1
  (1..s.length - 2).each { |i| s[i] += [a[i - 1], a[i]].max }
  a = s
end
puts a.max
