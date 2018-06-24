File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  swap = false
  r = []
  s.each do |i|
    if swap
      u, v = i.chomp(',').split('-').map(&:to_i)
      r[u], r[v] = r[v], r[u]
    elsif i == ':'
      swap = true
    else
      r << i
    end
  end
  puts r.join(' ')
end
