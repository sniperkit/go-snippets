File.open(ARGV[0]).each_line do |line|
  found = false
  s = line.split(',').map(&:to_i)
  s.uniq.each do |i|
    next if s.count(i) < s.length / 2
    break if s.count(i) == s.length / 2
    puts i
    found = true
    break
  end
  puts 'None' unless found
end
