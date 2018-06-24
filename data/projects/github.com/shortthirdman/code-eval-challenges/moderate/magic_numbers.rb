def magic?(a)
  d = []
  r = 0
  ns = []
  while a > 0
    r = a % 10
    return false if r == 0 || d[r]
    d[r] = true
    ns << r
    a /= 10
  end
  d = []
  r = 1
  ns.each do
    r = (r + ns[-r]) % ns.length
    r = r == 0 ? ns.length : r
    return false if d[r]
    d[r] = true
  end
  r == 1
end

maxb = 0
magic = []
File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  r = []
  while s[1] > maxb
    maxb += 1
    magic << maxb if magic?(maxb)
  end
  magic.each do |x|
    next if x < s[0]
    break if x > s[1]
    r << x
  end
  puts r.empty? ? -1 : r.join(' ')
end
