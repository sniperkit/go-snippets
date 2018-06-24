File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  s.delete_at(0) while !s.empty? && s.count(s[0]) == 1
  if s.length == 1
    puts
  else
    t = [s.delete_at(0)]
    s.each do |x|
      break if x == t[0]
      t << x
    end
    puts t.join(' ')
  end
end
