d = %w(zero hundred thousand million negative)
d0 = { 'zero' => 0, 'one' => 1, 'two' => 2, 'three' => 3, 'four' => 4,
       'five' => 5, 'six' => 6, 'seven' => 7, 'eight' => 8, 'nine' => 9 }
d1 = { 'ten' => 10, 'eleven' => 11, 'twelve' => 12, 'thirteen' => 13,
       'fourteen' => 14, 'fifteen' => 15, 'sixteen' => 16, 'seventeen' => 17,
       'eighteen' => 18, 'nineteen' => 19 }
d2 = { 'twenty' => 20, 'thirty' => 30, 'forty' => 40, 'fifty' => 50,
       'sixty' => 60, 'seventy' => 70, 'eighty' => 80, 'ninety' => 90 }

File.open(ARGV[0]).each_line do |line|
  h = q = r = 0
  t = line.chomp.split
  if t[0] == d[-1]
    n = -1
    t.shift
  else n = 1
  end
  t.each do |x|
    if x == d[0] then q = 0
    elsif d0[x] then h = d0[x]
    elsif d1[x] then q += d1[x]
    elsif d2[x] then q += d2[x]
    elsif x == d[1]
      q = h * 100
      h = 0
    elsif x == d[2]
      r += (h + q) * 1000
      h = q = 0
    elsif x == d[3]
      r = (h + q) * 1_000_000
      h = q = 0
    end
  end
  puts n * (r + q + h)
end
