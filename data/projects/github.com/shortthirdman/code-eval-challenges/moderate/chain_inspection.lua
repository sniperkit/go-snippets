for line in io.lines(arg[1]) do
  n, b, m, been = 0, "BEGIN", {}, {}
  for i in line:gmatch("[^;]+") do
    t, u = i:match("([^-]+)-([^-]+)")
    m[t], n = u, n + 1
  end
  for i = 1, n do
    b = m[b]
    if not b or been[b] then
      b = nil
      break
    end
    been[b] = true
  end
  print((b == "END") and "GOOD" or "BAD")
end
