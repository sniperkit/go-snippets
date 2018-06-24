File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(';')
  a = s[0].split(',').map(&:to_i)
  b = s[1].split(',').map(&:to_i)
  j = 0
  r = []
  a.each do |x|
    j += 1 while j < b.length && x > b[j]
    if j < b.length && x == b[j]
      r << x.to_s
      j += 1
    end
  end
  puts r.join(',')
end
