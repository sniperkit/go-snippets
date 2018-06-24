File.open(ARGV[0]).each_line do |line|
  line.chomp!
  if line.empty?
    puts 'NONE'
    next
  end
  l = line.length
  m = ''
  n = 1
  start = 0
  while n <= (l - start) / 2
    f = false
    (start...l - n).each do |i|
      s = line[i...i + n]
      next if s.strip.empty?
      next unless line[i + n..-1].include? s
      m = s
      start = i
      f = true
      n += 1
      break
    end
    break unless f
  end
  puts m.empty? ? 'NONE' : m
end
