Node = Struct.new(:value, :next)
def list_push(value, node)
  Node.new(value, node)
end

File.open(ARGV[0]).each_line do |line|
  n, m = line.split(',').map(&:to_i)

  tail = Node.new(n - 1, nil)
  list = tail
  (n - 2).step(0, -1) { |x| list = list_push(x, list) }
  tail.next = list
  list = tail

  s = []
  n.times do
    (m - 1).times { list = list.next }
    s << list.next.value
    list.next = list.next.next
  end
  puts s.join(' ')
end
