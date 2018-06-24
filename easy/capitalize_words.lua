for line in io.lines(arg[1]) do
    print((line:gsub("^%l", string.upper):gsub("%s%l", string.upper)))
end
