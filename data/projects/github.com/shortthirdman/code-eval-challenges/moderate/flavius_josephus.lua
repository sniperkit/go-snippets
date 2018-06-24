for line in io.lines(arg[1]) do
  n, m = line:match("(%d+),(%d+)")
  tail = {next = nil, value = n-1}
  list = tail
  for i = n-2, 0, -1 do list = {next = list, value = i} end
  tail.next = list
  list = tail
  for i = 1, n do
    for _ = 1, m-1 do list = list.next end
    if i > 1 then io.write(' ') end
    io.write(list.next.value)
    list.next = list.next.next
  end
  io.write('\n')
end
