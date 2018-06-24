h = { 'a' => 0, 'b' => 1, 'c' => 2, 'd' => 3, 'e' => 4,
      'f' => 5, 'g' => 6, 'h' => 7, 'i' => 8, 'j' => 9 }
File.open(ARGV[0]).each_line do |line|
  s = line.gsub(/[^\da-j]/, '').gsub(/[a-j]/, h)
  puts s.length > 0 ? s : 'NONE'
end
