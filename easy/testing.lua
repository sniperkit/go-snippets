local m = {"Done", "Low", "Medium", "High", "Critical"}

function prio(a)
  if a == 0 then
    return 1
  elseif a <= 2 then
    return 2
  elseif a <= 4 then
    return 3
  elseif a <= 6 then
    return 4
  else
    return 5
  end
end

for line in io.lines(arg[1]) do
  local sep = line:find("|")
  local r = 0
  for i = 1, sep - 2 do
    if line:sub(i, i) ~= line:sub(sep+1+i, sep+1+i) then
      r = r + 1
    end
  end
  print(m[prio(r)])
end
