def balanced(s, c)
  while c >= 0
    return c == 0 if s.empty?
    return balanced(s[2..-1], c) || balanced(s[2..-1], c + 1) if s[0..1] == ':('
    return balanced(s[2..-1], c) || balanced(s[2..-1], c - 1) if s[0..1] == ':)'
    c += 1 if s[0] == '('
    c -= 1 if s[0] == ')'
    s = s[1..-1]
  end
  false
end

File.open(ARGV[0]).each_line do |line|
  puts balanced(line.chomp, 0) ? 'YES' : 'NO'
end
