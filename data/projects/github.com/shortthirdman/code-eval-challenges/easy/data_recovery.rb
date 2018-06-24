File.open(ARGV[0]).each_line do |line|
  next unless line.include?(';')
  s = line.split(';')
  s0 = s[0].split
  s1 = s[1].split.map(&:to_i)
  (1..s0.length).each do |i|
    s1 << i unless s1.include?(i)
  end
  puts s0.sort_by.with_index { |_, i| s1[i] }.join(' ')
end
