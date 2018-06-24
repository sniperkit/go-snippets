def find_path(p, d, mx)
  r = []
  (0...@n).each { |x| r << x unless p.include?(x) }
  return [p + r, d + @dis[p[-1]][r.first]] if r.length == 1
  m = mx
  mp = []
  r.each do |x|
    e = d + @dis[p[-1]][x]
    next unless e < m
    di = find_path(p + [x], e, m)
    mp, m = di if di[1] < m
  end
  [mp, m]
end

def dist(a, b)
  x = a[1] - b[1]
  y = a[2] - b[2]
  Math.sqrt(x * x + y * y)
end

loc = []
File.open(ARGV[0]).each_line do |line|
  s = line.scan(/([0-9]+)[^\(]+\(([-0-9.]+), ([-0-9.]+)/)[0]
  loc << [s[0].to_i, s[1].to_f, s[2].to_f]
end

@dis = []
loc.each do |x|
  @dis << []
  loc.each { |y| @dis[-1] << dist(x, y) }
end
@n = @dis.length

p = find_path([0], 0, Float::MAX)
p[0].each { |x| puts x + 1 }
