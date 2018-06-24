morse = 'ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90'
File.open(ARGV[0]).each_line do |line|
  m = 1
  line.each_char do |i|
    case i
    when '.' then m *= 2
    when '-' then m = m * 2 + 1
    else
      if m == 1
        print ' '
      else
        print morse[m - 2] if m < 64
        m = 1
      end
    end
  end
  print morse[m - 2] if m > 1 && m < 64
  puts
end
