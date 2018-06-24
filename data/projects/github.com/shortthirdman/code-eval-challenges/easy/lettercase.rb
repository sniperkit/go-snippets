File.open(ARGV[0]).each_line do |line|
  l = line.count('a-z').to_f
  u = line.count('A-Z').to_f
  l, u = 100 * l / (l + u), 100 * u / (l + u) if l + u > 0
  printf("lowercase: %.2f uppercase: %.2f\n", l, u)
end
