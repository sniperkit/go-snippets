# permutation
class String
  def permutation
    chars.to_a.permutation.map(&:join)
  end
end

File.open(ARGV[0]).each_line do |line|
  puts line.chomp.permutation.sort.join(',')
end
