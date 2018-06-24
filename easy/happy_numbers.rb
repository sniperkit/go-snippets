def happy(a)
  ret = 0
  while a > 0
    b = a % 10
    ret += b * b
    a /= 10
  end
  ret
end

File.open(ARGV[0]).each_line do |line|
  a = line.to_i
  b = [a]
  127.times do
    break if a <= 1
    a = happy(a)
    if b.include?(a)
      a = 0
      break
    end
    b << a
  end
  puts a == 1 ? 1 : 0
end
