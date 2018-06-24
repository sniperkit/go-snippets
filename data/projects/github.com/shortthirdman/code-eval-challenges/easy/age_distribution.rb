cat = ['This program is for humans',
       "Still in Mama's arms",
       'Preschool Maniac',
       'Elementary school',
       'Middle school',
       'High school',
       'College',
       'Working for the man',
       'The Golden Years']
age = [0, 3, 5, 12, 15, 19, 23, 66, 101]

File.open(ARGV[0]).each_line do |line|
  a = line.to_i
  c = 0
  c += 1 while c < 9 && a >= age[c]
  puts cat[c % 9]
end
