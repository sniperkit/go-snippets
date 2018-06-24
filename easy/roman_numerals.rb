ronum = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]
rostr = %w(M CM D CD C XC L XL X IX V IV I)

File.open(ARGV[0]).each_line do |line|
  i = 0
  num = line.to_i
  while num > 0
    if num >= ronum[i]
      print rostr[i]
      num -= ronum[i]
    else
      i += 1
    end
  end
  puts
end
