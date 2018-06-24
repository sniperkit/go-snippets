File.open(ARGV[0]).each_line do |line|
  s = line.split(';')
  next if s.empty?
  m = {}
  s[1].split(',').map(&:to_i).each do |i|
    if m[i]
      puts i
      break
    end
    m[i] = true
  end
end
