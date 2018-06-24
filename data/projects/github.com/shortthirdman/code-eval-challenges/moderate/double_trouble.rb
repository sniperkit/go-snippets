File.open(ARGV[0]).each_line do |s|
  r = 1
  n = s.length / 2
  (0..n).each do |i|
    if (s[i] == 'A' && s[i + n] == 'B') || (s[i] == 'B' && s[i + n] == 'A')
      r = 0
      break
    end
    r *= 2 if s[i] == '*' && s[i + n] == '*'
  end
  puts r
end
