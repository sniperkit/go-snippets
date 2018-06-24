@moy = { 'Jan' => 1, 'Feb' => 2, 'Mar' => 3, 'Apr' => 4, 'May' => 5,
         'Jun' => 6, 'Jul' => 7, 'Aug' => 8, 'Sep' => 9, 'Oct' => 10,
         'Nov' => 11, 'Dec' => 12 }
def month(s)
  year = s[4, 4].to_i
  12 * (year - 1990) + @moy[s[0, 3]] - 1
end

File.open(ARGV[0]).each_line do |line|
  work = {}
  line.chomp.split('; ').each do |i|
    t0 = month(i[0, 8])
    t1 = month(i[9, 8])
    (t0..t1).each do |j|
      work[j] = true
    end
  end
  puts work.length / 12
end
