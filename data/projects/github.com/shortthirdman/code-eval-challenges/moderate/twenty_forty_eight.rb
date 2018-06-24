def slide(a)
  last_num = last_num_id = last_zero = -1
  (0...a.length).each do |i|
    if a[i] == 0
      last_zero = i if last_zero == -1
      next
    elsif last_num == a[i]
      a[last_num_id] = 2 * last_num
      last_num = -1
      a[i] = 0
      last_zero = i if last_zero == -1
      next
    end
    last_num = a[i]
    if last_zero == -1
      last_num_id = i
      next
    end
    last_num_id = last_zero
    a[last_zero] = a[i]
    a[i] = 0
    last_zero += 1
  end
  a
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split('; ')
  n = s[1].to_i
  t = []
  s[2].split('|').each do |x|
    v = []
    x.split.each { |y| v << y.to_i }
    t << v
  end
  (0...n).each do |x|
    case s[0]
    when 'LEFT'
      t[x] = slide(t[x])
    when 'RIGHT'
      t[x] = slide(t[x].reverse).reverse
    when 'UP'
      v = []
      (0...n).each { |y| v << t[y][x] }
      v = slide(v)
      (0...n).each { |y| t[y][x] = v[y] }
    when 'DOWN'
      v = []
      (0...n).each { |y| v << t[n - y - 1][x] }
      v = slide(v)
      (0...n).each { |y| t[n - y - 1][x] = v[y] }
    end
  end
  (0...n).each { |x| t[x] = t[x].join(' ') }
  puts t.join('|')
end
