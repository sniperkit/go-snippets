m = %w(Done Low Medium High Critical)

def prio(a)
  case a
  when 0 then 0
  when 1..2 then 1
  when 3..4 then 2
  when 5..6 then 3
  else 4
  end
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(' | ')
  r = 0
  (0...s[0].length).each do |i|
    r += 1 if s[0][i] != s[1][i]
  end
  puts m[prio(r)]
end
