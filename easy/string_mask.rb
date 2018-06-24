File.open(ARGV[0]).each_line do |line|
  s = line.split
  (0..s[0].length).each do |i|
    print s[1][i] == '1' ? s[0][i].upcase : s[0][i]
  end
  puts
end
