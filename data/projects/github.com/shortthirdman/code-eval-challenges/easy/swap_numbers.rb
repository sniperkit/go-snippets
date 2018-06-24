File.open(ARGV[0]).each_line do |line|
  puts line.split.map { |x| x[-1] + x[1..-2] + x[0] }.join(' ')
end
