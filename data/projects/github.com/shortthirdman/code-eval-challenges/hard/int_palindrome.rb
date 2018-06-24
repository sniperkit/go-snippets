def palindrome?(a)
  a == a.to_s.reverse.to_i
end

File.open(ARGV[0]).each_line do |line|
  l, r = line.split.map(&:to_i)
  n = 0
  (l..r).each do |i|
    prev = -1
    (i..r).each do |j|
      if prev >= 0
        p = prev
        p += 1 if palindrome?(j)
      else
        p = 0
        (i..j).each do |k|
          p += 1 if palindrome?(k)
        end
      end
      n += 1 if p.even?
      prev = p
    end
  end
  puts n
end
