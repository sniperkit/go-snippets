def g2d(a)
  a ^= a >> 4
  a ^= a >> 2
  a ^= a >> 1
  a
end

File.open(ARGV[0]).each_line do |line|
  puts line.split('|').map { |x| g2d(x.to_i(2)) }.join(' | ')
end
