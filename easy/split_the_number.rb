File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  next if s.empty?
  c = r = v = 0
  o = 1
  s[1].each_char do |i|
    if i =~ /[a-z]/
      v = v * 10 + s[0][c].to_i
      c += 1
    elsif i == '+'
      r += v * o
      o = 1
      v = 0
    elsif i == '-'
      r += v * o
      o = -1
      v = 0
    end
  end
  puts r + v * o
end
