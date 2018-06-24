regex1 = "^[%w.+-=]+@%w+[%w.]*%w*%.%w%w%w?%w?$"
regex2 = "^\"[%w@.+-=]+\"@%w+[%w.]*%w*%.%w%w%w?%w?$"

for line in io.lines(arg[1]) do
  print((line:find(regex1) or line:find(regex2)) and "true" or "false")
end
