function div(a, b) return math.floor(a/b) end

function binary(sl, a)
  if a == 0 then
    return string.rep("0", sl)
  else
    local b, c, ret = 0, a, ""
    while a > 0 do
      b, a = b * 2 + a % 2, div(a, 2)
    end
    while c > 0 do
      ret = ret .. tostring(b % 2)
      b, c = div(b, 2), div(c, 2)
    end
    return string.rep("0", sl-#ret) .. ret
  end
end

function genkey()
  local sl, c = 1, -1
  return function()
           while true do
             c = c + 1
             if c < 2^sl-1 then
               return binary(sl, c)
             end
             sl, c = sl + 1, -1
           end
  end
end

inp, i = "", 0
for line in io.lines(arg[1]) do
  inp = inp .. line
  if inp:sub(-4) == "1000" then
    if i == 0 then
      nextKey, trans, i = genkey(), {}, 1
      while inp:sub(i, i) ~= "0" and inp:sub(i, i) ~= "1" do
        trans[nextKey()] = inp:sub(i, i)
        i = i + 1
      end
      head, curr, l = inp:sub(i, i), "", 0
    end
    while i < #line do
      i = i + 1
      if #head < 3 then
        head = head .. inp:sub(i, i)
        if head == "000" then
          print()
          inp, i = "", 0
          break
        end
      else
        if l == 0 then
          l = tonumber(head, 2)
        end
        curr = curr .. inp:sub(i, i)
        if #curr == l then
          if curr == string.rep("1", l) then
            head, curr, l = "", "", 0
          else
            io.write(trans[curr])
            curr = ""
          end
        end
      end
    end
  end
end
