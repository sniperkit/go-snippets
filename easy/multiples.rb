File.open(ARGV[0]).each_line do |line|
  s = line.split(',').map(&:to_i)
  c = s[0] - (s[0] * 2**-Math.log2(s[1])).floor * 2**Math.log2(s[1])
  puts c > 0 ? s[0] - c.to_i + s[1] : s[0]
end
