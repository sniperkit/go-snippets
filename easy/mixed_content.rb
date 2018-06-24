File.open(ARGV[0]).each_line do |line|
  v = []
  w = []
  line.chomp.split(',').each do |x|
    if x =~ /^\d+$/
      v << x
    else
      w << x
    end
  end
  if w.length > 0
    print w.join(',')
    print '|' if v.length > 0
  end
  print v.join(',') if v.length > 0
  puts
end
