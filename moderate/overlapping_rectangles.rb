def within_rect(rect, pt)
  pt[0] >= rect[0] && pt[0] <= rect[2] && pt[1] >= rect[3] && pt[1] <= rect[1]
end

File.open(ARGV[0]).each_line do |line|
  s = line.split(',').map(&:to_i)
  r1 = s.first(4)
  r2 = s.last(4)
  p1 = [[r1[0], r1[1]], [r1[0], r1[3]], [r1[2], r1[1]], [r1[2], r1[3]]]
  p2 = [[r2[0], r2[1]], [r2[0], r2[3]], [r2[2], r2[1]], [r2[2], r2[3]]]

  f = false
  p1.each do |x|
    if within_rect(r2, x)
      f = true
      break
    end
  end
  unless f
    p2.each do |x|
      if within_rect(r1, x)
        f = true
        break
      end
    end
  end
  puts f ? 'True' : 'False'
end
