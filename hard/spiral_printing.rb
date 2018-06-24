File.open(ARGV[0]).each_line do |line|
  s = line.split(';')
  row = s[0].to_i
  col = s[1].to_i
  t = s[2].split
  i = tn = tw = 0
  j = 1
  te = col - 1
  ts = row - 1
  print t[0]
  loop do
    while j <= te
      print ' ' + t[i * col + j]
      j += 1
    end
    j -= 1
    i += 1
    tn += 1
    break if i > ts
    while i <= ts
      print ' ' + t[i * col + j]
      i += 1
    end
    i -= 1
    j -= 1
    te -= 1
    break if j < tw
    while j >= tw
      print ' ' + t[i * col + j]
      j -= 1
    end
    j += 1
    i -= 1
    ts -= 1
    break if i < tn
    while i >= tn
      print ' ' + t[i * col + j]
      i -= 1
    end
    i += 1
    j += 1
    tw += 1
    break if j > te
  end
  puts
end
