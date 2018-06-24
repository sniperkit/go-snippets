for line in io.lines(arg[1]) do
  local r = line:match(", (.*)")
  print((line:match("^(.*), "):gsub("[" .. r .. "]", ""):gsub("^%s+", ""):gsub("%s+$", "")))
end
