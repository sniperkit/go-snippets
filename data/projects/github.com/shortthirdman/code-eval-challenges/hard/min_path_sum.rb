n = i = 0
b = []
File.open(ARGV[0]).each_line do |line|
  if n == 0
    n = line.to_i
    i = n
    next
  end
  s = line.split(',')
  if i == n
    b = Array.new(n)
    b[0] = s[0].to_i
    (1...n).each do |j|
      b[j] = b[j - 1] + s[j].to_i
    end
  else
    b[0] += s[0].to_i
    (1...n).each do |j|
      b[j] = [b[j - 1], b[j]].min + s[j].to_i
    end
  end
  i -= 1
  if i == 0
    puts b[-1]
    n = 0
  end
end
