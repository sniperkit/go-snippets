lol = []
n = 0
File.open(ARGV[0]).each_line do |line|
  if n == 0
    n = line.to_i
    next
  end
  line.chomp!
  (0...n).each do |i|
    if i == lol.length
      lol[i] = line
      break
    end
    lol[i], line = line, lol[i] if line.length > lol[i].length
  end
end
puts lol.join("\n")
