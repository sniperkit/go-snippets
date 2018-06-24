File.open(ARGV[0]).each_line do |line|
  s = line.delete('|').split.map(&:to_i)
  0.upto(s.length / 2 - 1) do |i|
    print ' ' if i > 0
    print s[i] * s[s.length / 2 + i]
  end
  puts
end
