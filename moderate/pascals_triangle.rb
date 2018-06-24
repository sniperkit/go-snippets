File.open(ARGV[0]).each_line do |line|
  print 1
  (1...line.to_i).each do |i|
    print ' 1'
    r = 1
    i.times do |j|
      r = (r * (i - j)) / (j + 1)
      print ' ' + r.to_s
    end
  end
  puts
end
