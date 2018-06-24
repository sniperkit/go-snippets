def decode(s)
  loop do
    return 1 if s.length < 2
    # should never start with zero: let's ignore if it does
    if s[0].to_i > 2 || s[0] == '0'
      s = s[1..-1]
    elsif s[1] == '0' || (s[0] == '2' && s[1].to_i > 6)
      s = s[2..-1]
    else
      return s.length == 2 ? 2 : decode(s[1..-1]) + decode(s[2..-1])
    end
  end
end

File.open(ARGV[0]).each_line do |line|
  puts decode(line.chomp)
end
