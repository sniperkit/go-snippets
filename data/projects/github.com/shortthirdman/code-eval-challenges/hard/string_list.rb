File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(',')
  n = s[0].to_i
  t = s[1].split('').uniq * n
  puts t.sort.permutation(n).to_a.uniq.map(&:join).join(',')
end
