b = [0, 1, 2, 1, 2]
File.open(ARGV[0]).each_line do |line|
  a = line.to_i
  puts a / 5 + b[a % 5]
end
