def btr(c, x, y, i, j)
  return '' if i < 1 || j < 1
  return btr(c, x, y, i - 1, j - 1) + x[i - 1] if x[i - 1] == y[j - 1]
  return btr(c, x, y, i, j - 1) if c[i][j - 1] > c[i - 1][j]
  btr(c, x, y, i - 1, j)
end

File.open(ARGV[0]).each_line do |line|
  line.chomp!
  next if line.empty?
  x, y = line.split(';')
  c = Array.new(x.length + 1) { |_| Array.new(y.length + 1) { |_| 0 } }
  x.each_char.with_index do |i, ix|
    y.each_char.with_index do |j, jx|
      if i == j
        c[ix + 1][jx + 1] = c[ix][jx] + 1
      else
        c[ix + 1][jx + 1] = [c[ix + 1][jx], c[ix][jx + 1]].max
      end
    end
  end
  puts btr(c, x, y, x.length, y.length)
end
