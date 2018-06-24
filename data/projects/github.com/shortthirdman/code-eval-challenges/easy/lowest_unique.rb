File.open(ARGV[0]).each_line do |line|
  t = line.split.map(&:to_i)
  s = t.sort
  s.delete(s[0]) while s.length > 1 && s[0] == s[1]
  puts s.length == 0 ? 0 : t.find_index(s[0]) + 1
end
