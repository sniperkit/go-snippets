for line in io.lines(arg[1]) do
  local lo, hi = 0, nil
  for i in line:gmatch("%S+") do
    if not hi then
      hi = tonumber(i)
    elseif i == "Lower" then
      hi = math.ceil((lo+hi)/2)-1
    elseif i == "Higher" then
      lo = math.ceil((lo+hi)/2)+1
    end
  end
  print(math.ceil((lo+hi)/2))
end
