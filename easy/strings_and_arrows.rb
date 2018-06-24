File.open(ARGV[0]).each_line do |line|
  puts line.scan(/(?=<--<<)/).length + line.scan(/(?=>>-->)/).length
end
