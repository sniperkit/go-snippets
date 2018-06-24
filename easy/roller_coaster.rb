File.open(ARGV[0]).each_line do |line|
  up = true
  line.chomp.each_char do |i|
    if i =~ /[[:alpha:]]/
      print up ? i.upcase : i.downcase
      up = !up
    else
      print i
    end
  end
  puts
end
