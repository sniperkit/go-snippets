for line in io.lines(arg[1]) do
  repeat
    local l = #line
    line = line:gsub("(.)%1", "%1")
  until l == #line
  print(line)
end
