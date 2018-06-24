File.open(ARGV[0]).each_line do |line|
  s = line.tr('()', '').split
  p = s[0].split(',').map(&:to_f)
  q = s[1].split(',').map(&:to_f)
  if p.length < 2 || q.length < 2
    puts 0
    next
  end
  if p[0] != 0
    (p.length - 1).step(0, -1) do |i|
      p[i] -= p[0]
    end
  end
  if q[0] != 0
    (q.length - 1).step(0, -1) do |i|
      q[i] -= q[0]
    end
  end
  ins = 0
  stl = p[-1]
  avl = q[-1]
  (0...q.length).each do |i|
    q[i] = q[i] * stl / avl
  end
  i = j = 0
  while i < p.length && j < q.length
    if p[i] == q[j]
      ins += 1
      i += 1
      j += 1
    elsif p[i] < q[j]
      i += 1
    else
      j += 1
    end
  end
  puts p.length + q.length - ins - 1
end
