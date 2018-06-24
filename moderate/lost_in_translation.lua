codel = [[nug
rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp
tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi
de kr kd eoya kw aej icfkici re zjkr]]
decol = [[bjv
the public is amazed by the quickness of the juggler
we think that our language is impossible to understand
so it is okay if you decided to quit]]

m, misscod, missdec = {}, {}, {}
m[" "] = " "
for i = string.byte("a"), string.byte("z") do
  ix = codel:find(string.char(i))
  if ix then
    m[string.char(i)] = decol:sub(ix, ix)
  else
    misscod[#misscod+1] = string.char(i)
  end
  ix = decol:find(string.char(i))
  if not ix then
    missdec[#missdec+1] = string.char(i)
  end
end
for ix, i in ipairs(misscod) do
  m[i] = missdec[ix]
end

for line in io.lines(arg[1]) do
  local ret = ""
  for i in line:gmatch("[a-z ]") do
    ret = ret .. m[i]
  end
  print(ret)
end
