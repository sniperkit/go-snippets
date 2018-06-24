File.open(ARGV[0]).each_line do |line|
  puts line.split.sort { |x, y| x.to_f <=> y.to_f }.join(' ')
end
