def char(a)
  (a % 26 + 'A'.ord).chr
end

File.open(ARGV[0]).each_line do |a|
  a = a.to_i - 1
  d = char(a)
  if a >= 26
    a = a / 26 - 1
    d = char(a) + d
    if a >= 26
      a = a / 26 - 1
      d = char(a) + d
    end
  end
  puts d
end
