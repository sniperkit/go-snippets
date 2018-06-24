File.open(ARGV[0]).each_line do |line|
  p = 1
  (1...line.chomp.length).each do |i|
    p = i + 1 if line[i] != line[i - p]
  end
  puts p
end
