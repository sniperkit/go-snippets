File.open(ARGV[0]).each_line do |line|
  puts line.gsub(/(.)\1+/, '\\1')
end
