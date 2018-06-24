# rotate!
class String
  def rotate!
    self << self[0]
    self[0] = ''
  end
end

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(',')
  r = false
  s[0].length.times do
    if s[0] == s[1]
      r = true
      break
    end
    s[1].rotate!
  end
  puts r ? 'True' : 'False'
end
