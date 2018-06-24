File.open(ARGV[0]).each_line do |line|
  line.chomp!
  line.each_char do |i|
    if line.count(i) == 1
      print i
      break
    end
  end
  puts
end
