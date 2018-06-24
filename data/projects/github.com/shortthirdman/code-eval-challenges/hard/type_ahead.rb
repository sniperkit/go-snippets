TEXT = '
Mary had a little lamb its fleece was white as snow;
And everywhere that Mary went, the lamb was sure to go.
It followed her to school one day, which was against the rule;
It made the children laugh and play, to see a lamb at school.
And so the teacher turned it out, but still it lingered near,
And waited patiently about till Mary did appear.
"Why does the lamb love Mary so?" the eager children cry;
"Why, Mary loves the lamb, you know" the teacher did reply."
'.freeze

TXL = TEXT.delete('.,;?"').split

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(',')
  n = s[0].to_i
  t = s[1].split
  tnl = []
  (0..TXL.length - n).each { |x| tnl << TXL[x...x + n] }
  r = []
  tnl.each do |x|
    r << x if x.take(t.length) == t
  end
  r.sort! do |x, y|
    if r.count(x) == r.count(y)
      x.join(' ') <=> y.join(' ')
    else r.count(y) <=> r.count(x)
    end
  end
  u = []
  r.each_with_index do |i, ix|
    if ix == 0 || i != r[ix - 1]
      u << i[t.length..-1].join(' ') + format(',%.3f',
                                              r.count(i).to_f / r.length)
    end
  end
  puts u.join(';')
end
