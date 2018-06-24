File.open(ARGV[0]).each_line do |line|
  s = line.split.map(&:to_i)
  f = false
  s.permutation.each do |i|
    (0..2).each do |o1|
      r1 = case o1
           when 0 then i[0] + i[1]
           when 1 then i[0] - i[1]
           else i[0] * i[1]
           end
      (0..2).each do |o2|
        r2 = case o2
             when 0 then r1 + i[2]
             when 1 then r1 - i[2]
             else r1 * i[2]
             end
        (0..2).each do |o3|
          r3 = case o3
               when 0 then r2 + i[3]
               when 1 then r2 - i[3]
               else r2 * i[3]
               end
          (0..2).each do |o4|
            r4 = case o4
                 when 0 then r3 + i[4]
                 when 1 then r3 - i[4]
                 else r3 * i[4]
                 end
            if r4 == 42
              f = true
              break
            end
          end
          break if f
        end
        break if f
      end
      break if f
    end
    break if f
  end
  puts f ? 'YES' : 'NO'
end
