# min-heap
class Heap
  def initialize(len = 0)
    @nodes = (0...len).to_a
  end

  def empty
    @nodes.count == 0
  end

  def has(node)
    @nodes.include? node
  end

  def pop_min
    ret = @nodes[0]
    max = @nodes.pop
    if @nodes.count > 0
      @nodes[0] = max
      bubble_down max, 0
    end
    ret
  end

  def del_node(node)
    return false unless @nodes.include? node
    max = @nodes.pop
    return if node == max
    offset = @nodes.index(node)
    @nodes[offset] = max
    bubble_down max, offset
  end

  def push(node)
    offset = @nodes.count
    @nodes << node
    bubble_up node, offset
  end

  def bubble_up(node, offset)
    parent = (offset - 1) / 2
    while offset > 0 && @nodes[parent] > node
      @nodes[parent], @nodes[offset] = @nodes[offset], @nodes[parent]
      offset = parent
      parent = (offset - 1) / 2
    end
  end

  def bubble_down(node, offset)
    while offset < @nodes.count / 2
      child1 = offset * 2 + 1
      child2 = child1 + 1
      min = @nodes[child2] && @nodes[child1] > @nodes[child2] ? child2 : child1
      break if node <= @nodes[min]
      @nodes[offset], @nodes[min] = @nodes[min], @nodes[offset]
      offset = min
    end
  end
end

File.open(ARGV[0]).each_line do |line|
  s = line.split(',').map(&:to_i)
  k = s[1]
  m = [s[2]]
  (k - 1).times do
    m << (s[3] * m[-1] + s[4]) % s[5]
  end

  h = Heap.new(k + 1)
  m[-k..-1].each do |i|
    h.del_node i
  end

  while m.length + 1 < s[0]
    n = h.pop_min
    h.push(m[-k]) unless h.has(m[-k]) || (m[-k + 1..-1].include? m[-k])
    m << n
  end
  puts h.pop_min
end
