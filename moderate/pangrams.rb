File.open(ARGV[0]).each_line do |line|
  f = false
  line = line.downcase
  ('a'..'z').each do |x|
    unless line.include?(x)
      print x
      f = true
    end
  end
  puts f ? '' : 'NULL'
end
