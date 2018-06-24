for line in io.lines(arg[1]) do
    print(math.fmod(tonumber(line) + 1, 2))
end
