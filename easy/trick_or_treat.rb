p = /Vampires: (\d+), Zombies: (\d+), Witches: (\d+), Houses: (\d+)/
File.open(ARGV[0]).each_line do |line|
  s = line.scan(p)[0].map(&:to_i)
  if s[0] + s[1] + s[2] == 0
    puts 0
  else
    puts s[3] * (s[0] * 3 + s[1] * 4 + s[2] * 5) / (s[0] + s[1] + s[2])
  end
end
