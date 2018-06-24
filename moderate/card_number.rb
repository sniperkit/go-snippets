File.open(ARGV[0]).each_line do |line|
  s = line.delete(" \n").split('').map(&:to_i)
  (s.length - 2).step(0, -2) do |i|
    s[i] *= 2
    s[i] = s[i] % 10 + 1 if s[i] > 9
  end
  puts s.inject(:+) % 10 == 0 ? 1 : 0
end
