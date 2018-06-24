for line in io.lines(arg[1]) do
  local a, num, ar, rr = 0, 0, 0, 0
  for i in line:gmatch("%d[MDCLXVI]") do
    a, b = tonumber(i:sub(1, 1)), i:sub(2, 2)
    if b == 'M' then r = 1000
    elseif b == 'D' then r = 500
    elseif b == 'C' then r = 100
    elseif b == 'L' then r = 50
    elseif b == 'X' then r = 10
    elseif b == 'V' then r = 5
    else r = 1
    end
    num = (r > rr) and num - ar * rr or num + ar * rr
    ar, rr = a, r
  end
  print(num + ar * rr)
end
