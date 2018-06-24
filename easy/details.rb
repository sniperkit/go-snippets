File.open(ARGV[0]).each_line do |line|
  puts line.scan(/X\.*Y/).map(&:length).min - 2
end
