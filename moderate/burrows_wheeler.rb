File.open(ARGV[0]).each_line do |line|
  s = line.chomp.chop.split('').map(&:ord)
  k = 0
  s.each_index do |ix|
    k = ix if s[ix] == '$'.ord
    s[ix] = (s[ix] << 16) + ix
  end
  s.sort!
  s.length.times do
    print (s[k] >> 16).chr
    k = s[k] % (1 << 16)
  end
  puts
end
