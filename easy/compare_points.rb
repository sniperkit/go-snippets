File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  if s[0] == s[2]
    if s[1] == s[3]
      puts 'here'
    else
      puts s[1] < s[3] ? 'N' : 'S'
    end
  elsif s[1] == s[3]
    puts s[0] < s[2] ? 'E' : 'W'
  elsif s[0] < s[2]
    puts s[1] < s[3] ? 'NE' : 'SE'
  else
    puts s[1] < s[3] ? 'NW' : 'SW'
  end
end
