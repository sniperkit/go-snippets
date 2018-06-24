File.open(ARGV[0]).each_line do |line|
  s = line.split(/[: ]/).map(&:to_i)
  t = ((s[0] - s[3]) * 3600 + (s[1] - s[4]) * 60 + s[2] - s[5]).abs
  printf("%02d:%02d:%02d\n", t / 3600, (t / 60) % 60, t % 60)
end
