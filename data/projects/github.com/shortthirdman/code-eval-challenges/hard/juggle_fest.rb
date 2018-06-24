c = []
j = []
sol = 1970

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split
  if s[0] == 'C'
    c << [s[1], s[2][2..-1].to_i, s[3][2..-1].to_i, s[4][2..-1].to_i, []]
  elsif s[0] == 'J'
    p = s[5].split(',').map { |x| x[1..-1].to_i }
    j << [s[1], s[2][2..-1].to_i, s[3][2..-1].to_i, s[4][2..-1].to_i, p]
  end
end
f = j.length / c.length

v = []
u = (0...j.length).to_a
while u.length > 0
  # each juggler applies to top contest
  # that hasn't rejected them yet
  cur = u.shift
  jug = j[cur]
  if jug[4].length == 0
    v << jug
    next
  end
  cha = c[jug[4].shift]
  fit = jug[1] * cha[1] + jug[2] * cha[2] + jug[3] * cha[3]
  # each contest available only for the best f jugglers
  cha[4] << [cur, fit]
  cha[4].sort! { |x, y| y[1] <=> x[1] }
  u << cha[4].pop[0] while cha[4].length > f && cha[4][-1][1] < cha[4][f - 1][1]
end
v.each do |x|
  best = cont = 0
  # if not among juggler's favorite contests, go for best fit
  c.each_with_index do |y, yx|
    next if y.length > f
    fit = x[1] * y[1] + x[2] * y[2] + x[3] * y[3]
    if fit > best
      best = fit
      cont = yx
    end
  end
  c[cont][4] << [x, best]
end

sum = 0
c[sol][4].each do |x|
  sum += x[0]
end
puts sum
