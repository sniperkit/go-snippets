units = ['PENNY', 'NICKEL', 'DIME', 'QUARTER', 'HALF DOLLAR', 'ONE'] +
        ['TWO', 'FIVE', 'TEN', 'TWENTY', 'FIFTY', 'ONE HUNDRED']
value = [1, 5, 10, 25, 50, 100, 200, 500, 1000, 2000, 5000, 10_000]

def centi(a)
  return a.to_i * 100 unless a.include? '.'
  a1, a2 = a.split('.').map(&:to_i)
  a1 * 100 + a2
end

File.open(ARGV[0]).each_line do |line|
  p, c = line.split(';').map { |i| centi(i) }
  if p > c
    puts 'ERROR'
    next
  elsif p == c
    puts 'ZERO'
    next
  end
  first = true
  (units.length - 1).downto(0) do |i|
    while p + value[i] <= c
      if first then first = false
      else print ','
      end
      print units[i]
      p += value[i]
    end
  end
  puts
end
