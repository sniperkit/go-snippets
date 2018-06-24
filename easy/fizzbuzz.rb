File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  (1..s[2]).each do |i|
    print ' ' if i > 1
    if i % s[0] > 0 && i % s[1] > 0
      print i
    else
      print 'F' if i % s[0] == 0
      print 'B' if i % s[1] == 0
    end
  end
  puts
end
