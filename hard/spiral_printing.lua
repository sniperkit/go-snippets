for line in io.lines(arg[1]) do
  row, col, line = line:match("(%d+);(%d+);(.*)")
  t = {}
  for i in string.gmatch(line, "%S+") do
    t[#t + 1] = i
  end
  local i, j = 0, 1
  local tn, te, ts, tw = 0, col - 1, row - 1, 0
  io.write(t[1])
  repeat
    while j <= te do
      io.write(" " .. t[i*col + j + 1])
      j = j + 1
    end
    j, i, tn = j-1, i+1, tn+1
    if i > ts then break end
    while i <= ts do
      io.write(" " .. t[i*col + j + 1])
      i = i + 1
    end
    i, j, te = i-1, j-1, te-1
    if j < tw then break end
    while j >= tw do
      io.write(" " .. t[i*col + j + 1])
      j = j - 1
    end
    j, i, ts = j+1, i-1, ts-1
    if i < tn then break end
    while i >= tn do
      io.write(" " .. t[i*col + j + 1])
      i = i - 1
    end
    i, j, tw = i+1, j+1, tw+1
  until j > te
  print()
end
