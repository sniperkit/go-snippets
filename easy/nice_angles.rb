File.open(ARGV[0]).each_line do |line|
  a = line.to_f
  b = a.floor
  print b.to_s + '.'
  a = (a - b) * 60
  b = a.floor
  printf("%02d'", b)
  a = (a - b) * 60
  b = a.floor
  printf("%02d\"\n", b)
end
