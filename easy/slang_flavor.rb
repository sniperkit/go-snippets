phrases = [', yeah!',
           ', this is crazy, I tell ya.',
           ', can U believe this?',
           ', eh?',
           ', aw yea.',
           ', yo.',
           '? No way!',
           '. Awesome!']
i = 0
l = false
File.open(ARGV[0]).each_line do |line|
  s = line.split('').each do |x|
    if '.!?'.include? x
      if l
        print phrases[i]
        i = (i + 1) % 8
        l = false
      else
        print x
        l = true
      end
    else
      print x
    end
  end
  puts if s.last != "\n"
end
