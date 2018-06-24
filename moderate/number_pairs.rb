File.open(ARGV[0]).each_line do |line|
  s = line.split(';')
  t = s[0].split(',').map(&:to_i)
  u = s[1].to_i
  n = false
  i = 0
  j = t.length - 1
  while i < j && 2 * t[i] < u
    while j > i && t[i] + t[j] >= u
      if t[i] + t[j] == u
        if n
          print ';'
        else
          n = true
        end
        printf '%d,%d', t[i], t[j]
      end
      j -= 1
    end
    i += 1
  end
  print 'NULL' unless n
  puts
end
