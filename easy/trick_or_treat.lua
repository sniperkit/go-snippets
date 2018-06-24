function div(a, b) return math.floor(a/b) end

for line in io.lines(arg[1]) do
  local v, z, w, h = line:match("Vampires: (%d+), Zombies: (%d+), Witches: (%d+), Houses: (%d+)")
  print(v + z + w == 0 and 0 or div((v * 3 + z * 4 + w * 5) * h, v + z + w))
end
