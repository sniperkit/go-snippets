sq = {}
n = 0
(0..46_340).each { |i| sq[i * i] = true }

File.open(ARGV[0]).each_line do |line|
  if n == 0
    n = line.to_i
  else
    x = line.to_i
    num = 0
    l = {}
    bot = Math.sqrt(x / 2).to_i
    top = Math.sqrt(x).to_i
    (bot..top).each do |i|
      t = x - i * i
      l[i * i] = true
      num += 1 if sq[t] && !l[t]
    end
    puts num
  end
end
