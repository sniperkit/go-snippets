Team = Struct.new(:id, :countries)

File.open(ARGV[0]).each_line do |line|
  s = line.split(' | ').map(&:split)
  t = []
  r = []
  s.each_with_index do |i, ix|
    i.each do |j|
      jx = j.to_i
      f = false
      (0...t.length).each do |k|
        next if t[k].id != jx
        t[k].countries += format(',%d', ix + 1)
        f = true
        break
      end
      t << Team.new(jx, (ix + 1).to_s) unless f
    end
  end
  t.sort { |x, y| x.id <=> y.id }.each do |i|
    r << format('%d:%s;', i.id, i.countries)
  end
  puts r.join(' ')
end
