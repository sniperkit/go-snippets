File.open(ARGV[0]).each_line do |line|
  s = line.chomp
  if s.length > 55
    p = 40
    39.downto(0) do |i|
      if s[i] == ' '
        p = i
        break
      end
    end
    s = s[0...p] + '... <Read More>'
  end
  puts s
end
