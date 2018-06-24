T = %w(0 1 abc def ghi jkl mno pqrs tuv wxyz).map { |x| x.split('') }.freeze

File.open(ARGV[0]).each_line do |line|
  s = line.chomp.split('').map(&:to_i)
  puts T[s[0]].product(T[s[1]], T[s[2]], T[s[3]], T[s[4]], T[s[5]], T[s[6]]
                      ).map(&:join).join(',')
end
