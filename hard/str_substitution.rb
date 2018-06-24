File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(';')
  if s.length < 2
    puts
    next
  end
  t = s[1].split(',')
  0.step(t.length - 1, 2) do |i|
    if t[i + 1].nil? || t[i + 1].empty?
      t[i + 1] = 'x'
    else
      t[i + 1].gsub!(/[01]/, '0' => 'a', '1' => 'b')
    end
    s[0].gsub!(/#{t[i]}/, t[i + 1])
  end
  puts s[0].delete('x').gsub(/[ab]/, 'a' => '0', 'b' => '1')
end
