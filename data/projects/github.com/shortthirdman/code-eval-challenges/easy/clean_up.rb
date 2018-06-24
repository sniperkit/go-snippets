File.open(ARGV[0]).each_line do |line|
  puts line.chomp.downcase.gsub(/^[^a-z]+|[^a-z]+$/, '').gsub(/[^a-z]+/, ' ')
end
