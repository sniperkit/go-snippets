File.open(ARGV[0]).each_line do |line|
  puts line.split('').each { |x|
    if x =~ /[[:lower:]]/
      x.upcase!
    elsif x =~ /[[:upper:]]/
      x.downcase!
    end
  }.join
end
