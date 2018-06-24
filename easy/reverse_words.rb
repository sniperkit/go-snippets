File.open(ARGV[0]).each_line do |line|
  puts line.split.reverse.join(' ') if line !~ /^$/
end
