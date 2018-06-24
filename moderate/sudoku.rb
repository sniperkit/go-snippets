@sq = { 4 => 2, 9 => 3 }
def region(a, n)
  z = @sq[n]
  (a % n) / z + z * ((a / n) / z)
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(';')
  g = s[1].split(',').map(&:to_i)
  n = s[0].to_i
  valid = true
  row = Array.new(9, 0)
  col = Array.new(9, 0)
  reg = Array.new(9, 0)
  g.each_with_index do |i, ix|
    if row[ix / n] & 1 << i > 0
      valid = false
      break
    end
    row[ix / n] |= 1 << i
    if col[ix % n] & 1 << i > 0
      valid = false
      break
    end
    col[ix % n] |= 1 << i
    if reg[region(ix, n)] & 1 << i > 0
      valid = false
      break
    end
    reg[region(ix, n)] |= 1 << i
  end
  puts valid ? 'True' : 'False'
end
