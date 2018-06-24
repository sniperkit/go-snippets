File.open(ARGV[0]).each_line do |line|
  s = line.delete(" \n").split('').map(&:to_i)
  (s.length - 2).step(0, -2) do |i|
    s[i] *= 2
  end
  puts s.inject(:+) % 10 == 0 ? 'Real' : 'Fake'
end
