def push(stack, int)
  stack << int
end

def pop(stack)
  int = stack[-1]
  stack.delete_at(-1)
  int
end

File.open(ARGV[0]).each_line do |line|
  stack = []
  line.split.map(&:to_i).each { |i| push(stack, i) }

  out = []
  while stack.length > 0
    out << pop(stack)
    pop(stack)
  end
  puts out.join(' ')
end
