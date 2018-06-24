k, o = [[ !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz]], string.byte("0")
for line in io.lines(arg[1]) do
  ix, r, s = 0, line:match("([^;]+);([^;]+)")
  for i in s:gmatch(".") do
    p = (#k + k:find(i, nil, true) - 1 - (r:byte(ix % #r + 1) - o)) % #k + 1
    io.write(k:sub(p, p))
    ix = ix + 1
  end
  print()
end
