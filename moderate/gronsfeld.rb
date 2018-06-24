k = %q{ !"#$%&'()*+,-./0123456789:<=>?@} +
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split(';')
  ix = 0
  s[1].each_char do |i|
    print k[(k.length + k.index(i) - s[0][ix % s[0].length].to_i) % k.length]
    ix += 1
  end
  puts
end
