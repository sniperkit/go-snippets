for line in io.lines(arg[1]) do
    print((line:gsub("%u", string.lower)))
end
