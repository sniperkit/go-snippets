V = { '2' => 2, '3' => 3, '4' => 4, '5' => 5, '6' => 6, '7' => 7, '8' => 8,
      '9' => 9, '1' => 10, 'J' => 11, 'Q' => 12, 'K' => 13, 'A' => 14 }.freeze
def rank(c, t)
  r = V[c[0]]
  r += 13 if c[-1] == t
  r
end

def win(l, r, t)
  case rank(l, t) <=> rank(r, t)
  when -1 then r
  when 1 then l
  else l + ' ' + r
  end
end

File.open(ARGV[0]).each_line do |line|
  s = line.split
  puts win(s[0], s[1], s[3])
end
