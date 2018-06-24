morse = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"
for line in io.lines(arg[1]) do
  local m = 1
  for i in line:gmatch(".") do
    if i == '.' then
      m = m * 2
    elseif i == '-' then
      m = m * 2 + 1
    elseif m == 1 then
      io.write(" ")
    else
      if m < 64 then
        io.write(morse:sub(m-1, m-1))
      end
      m = 1
    end
  end
  print((m < 64) and morse:sub(m-1, m-1) or "")
end
