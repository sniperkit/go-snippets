# sum
class Array
  def sum
    inject(0, :+)
  end
end

File.open(ARGV[0]).each_line do |line|
  s = line.split(' | ').map(&:split)
  if s[0].map { |x| x.to_i(16) }.sum > s[1].map { |x| x.to_i(2) }.sum
    puts 'False'
  else puts 'True'
  end
end
