File.open(ARGV[0]).each_line do |line|
  n, *s = line.chomp.split.map(&:to_i)
  jolly = true
  if n > 1
    u = Array.new(n - 1, false)
    (1...n).each do |i|
      x = (s[i] - s[i - 1]).abs
      if x >= n || x == 0 || u[x - 1]
        jolly = false
        break
      end
      u[x - 1] = true
    end
  end
  puts jolly ? 'Jolly' : 'Not jolly'
end
