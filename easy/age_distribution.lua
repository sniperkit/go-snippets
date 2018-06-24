local category = {"This program is for humans", "Still in Mama's arms", "Preschool Maniac", "Elementary school", "Middle school", "High school", "College", "Working for the man", "The Golden Years"}
local age = {0, 3, 5, 12, 15, 19, 23, 66, 101}

for line in io.lines(arg[1]) do
    local m = 0
    while m < 9 and tonumber(line) >= age[m+1] do
        m = m + 1
    end
    print(category[(m%9)+1])
end
