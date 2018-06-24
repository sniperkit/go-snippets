@h = {}

def overlap(a, b)
  k = a + ';' + b
  return @h[k] if @h[k]
  if a.length < b.length
    o = overlap(b, a)
    return [o[0], o[2], o[1]]
  end
  maxfit = [0, 0, 0]
  (0..a.length - 1).each do |i|
    l = [a.length - i, b.length].min
    if a[i...(i + l)] == b[0...l]
      maxfit = [l, i, 0]
      break
    end
  end
  (0..b.length).each do |i|
    l = b.length - i
    break if maxfit[0] >= l
    if a[0...l] == b[i...(i + l)]
      maxfit = [l, 0, i]
      break
    end
  end
  @h[k] = maxfit
  maxfit
end

File.open(ARGV[0]).each_line do |line|
  f = line[0...-1].split(';')
  while f.length > 1
    maxoverlap = [0, 0, 0, 0, 0]
    f.each_with_index do |i, ix|
      f.each_with_index do |j, jx|
        next if ix == jx
        o = overlap(i, j)
        maxoverlap = [o[0], ix, jx, o[1], o[2]] if o[0] > maxoverlap[0]
      end
    end
    if maxoverlap[4] == 0
      f[maxoverlap[1]] = f[maxoverlap[1]] << f[maxoverlap[2]][maxoverlap[0]..-1]
    elsif maxoverlap[3] == 0
      f[maxoverlap[1]] = f[maxoverlap[2]] << f[maxoverlap[1]][maxoverlap[0]..-1]
    end
    f.delete_at(maxoverlap[2])
  end
  puts f[0]
end
