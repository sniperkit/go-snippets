units = {"PENNY", "NICKEL", "DIME", "QUARTER", "HALF DOLLAR", "ONE",
	 "TWO", "FIVE", "TEN", "TWENTY", "FIFTY", "ONE HUNDRED"}
value = {1, 5, 10, 25, 50, 100, 200, 500, 1000, 2000, 5000, 10000}

function centi(a)
  if a:find("%.") then
    local a1, a2 = a:match("(%d+)%.(%d+)")
    return a1 * 100 + a2
  else
    return a * 100
  end
end

for line in io.lines(arg[1]) do
  local p, c = line:match("([^;]+);([^;]+)")
  p, c = centi(p), centi(c)
  if p > c then
    print("ERROR")
  elseif p == c then
    print("ZERO")
  else
    local first = true
    for i = #units, 1, -1 do
      while p + value[i] <= c do
        if not first then io.write(",") else first = false end
        io.write(units[i])
        p = p + value[i]
      end
    end
    print()
  end
end
