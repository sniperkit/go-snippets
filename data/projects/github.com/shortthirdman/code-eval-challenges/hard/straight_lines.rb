Point = Struct.new(:x, :y)

def in_line(a, b, c)
  (a.x - b.x) * (a.y - c.y) == (a.x - c.x) * (a.y - b.y)
end

def line(ax, bx, t)
  t.each_with_index do |i, ix|
    return ix >= bx if ix != ax && ix != bx && in_line(i, t[ax], t[bx])
  end
  false
end

File.open(ARGV[0]).each_line do |line|
  s = line.split(' | ').map do |i|
    t = i.split.map(&:to_i)
    Point.new(t[0], t[1])
  end
  r = 0
  (0..s.length - 2).each do |i|
    (i + 1...s.length).each do |j|
      r += 1 if line(i, j, s)
    end
  end
  puts r
end
