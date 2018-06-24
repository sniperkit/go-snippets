parent = { 3 => 8, 8 => 30, 10 => 20, 20 => 8, 29 => 20, 30 => nil, 52 => 30 }

File.open(ARGV[0]).each_line do |line|
  n = line.split.map(&:to_i)
  found = false
  while n[0] && n[0] != n[1]
    i = n[1]
    while i
      if n[0] == i
        found = true
        break
      end
      i = parent[i]
    end
    break if found
    n[0] = parent[n[0]]
  end
  puts n[0]
end
