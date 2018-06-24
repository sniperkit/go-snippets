File.open(ARGV[0]).each_line do |line|
  s = line.split
  lo = 0
  hi = s[0].to_i
  (1...s.length).each do |x|
    m = lo + hi
    c = m % 2
    case s[x]
    when 'Lower' then hi = m / 2 + c - 1
    when 'Higher' then lo = m / 2 + c + 1
    else puts m / 2 + c
    end
  end
end
