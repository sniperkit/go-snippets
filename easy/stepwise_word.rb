File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  m = ''
  s.each { |i| m = i if i.length > m.length }
  puts m.split('').each_with_index.map { |i, ix| '*' * ix + i }.join(' ')
end
