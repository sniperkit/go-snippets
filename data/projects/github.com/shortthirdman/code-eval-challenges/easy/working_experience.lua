function div(a, b) return math.floor(a/b) end

moy = {Jan=1, Feb=2, Mar=3, Apr=4, May=5, Jun=6,
       Jul=7, Aug=8, Sep=9, Oct=10, Nov=11, Dec=12}
function month(s)
  return 12*(tonumber(s:sub(5, 8))-1990) + moy[s:sub(1, 3)] - 1
end

for line in io.lines(arg[1]) do
  local work, w = {}, 0
  for i in line:gsub("; ", ";"):gmatch("[^;]+") do
    for j = month(i:sub(1, 8)), month(i:sub(10, 17)) do
      work[j] = true
    end
  end
  for _, _ in pairs(work) do w = w + 1 end
  print(div(w, 12))
end
