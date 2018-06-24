File.open(ARGV[0]).each_line do |line|
  puts line.chars.map {|x|
    case x
    when 'a'..'f' then (x.ord + 20).chr
    when 'g'..'m' then (x.ord + 7).chr
    when 'n'..'t' then (x.ord - 7).chr
    when 'u'..'z' then (x.ord - 20).chr
    end
  }.join
end
