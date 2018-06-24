File.open(ARGV[0]).each_line do |line|
  s = line.split
  t = []
  l = Math.sqrt(s.length).to_i
  0.upto(l - 1) do |x|
    (l - 1).downto(0) do |y|
      t << s[y * l + x]
    end
  end
  puts t.join(' ')
end
