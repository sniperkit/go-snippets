for line in io.lines(arg[1]) do
  local l = select(2, line:gsub("[a-z]", "."))
  local u = select(2, line:gsub("[A-Z]", "."))
  if l+u > 0 then
    print(string.format("lowercase: %.2f uppercase: %.2f", 100*l/(l+u), 100*u/(l+u)))
  else
    print("lowercase: 0.00 uppercase: 0.00")
  end
end
