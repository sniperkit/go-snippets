File.open(ARGV[0]).each_line do |line|
  s = line.gsub(/[^\d ]/, '').split.map(&:to_i).sort
  (-s.length + 1..-1).each { |i| s[-i] -= s[-i - 1] }
  puts s.join(',')
end
