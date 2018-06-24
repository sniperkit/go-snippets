File.open(ARGV[0]).each_line do |line|
  b = 'BEGIN'
  s = line.chomp.split(';')
  m = {}
  been = {}
  s.each do |x|
    t = x.split('-')
    m[t[0]] = t[1]
  end
  s.length.times do
    if !m.key?(b) || been.key?(m[b])
      b = nil
      break
    end
    b = m[b]
    been[b] = true
  end
  puts b == 'END' ? 'GOOD' : 'BAD'
end
