File.open(ARGV[0]).each_line do |line|
  a = num = ar = rr = 0
  line.chomp.each_char do |x|
    if x =~ /\d/ then a = x.to_i
    else
      case x
      when 'M' then r = 1000
      when 'D' then r = 500
      when 'C' then r = 100
      when 'L' then r = 50
      when 'X' then r = 10
      when 'V' then r = 5
      when 'I' then r = 1
      end

      if r > rr then num -= ar * rr
      else num += ar * rr
      end
      ar = a
      rr = r
    end
  end
  puts num + ar * rr
end
