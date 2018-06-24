P = { '^' => 3, '*' => 2, '/' => 2, '+' => 1, '-' => 1 }.freeze
ERROR = -999_999

def operate(a, b, c)
  case b
  when '^'
    return a**c
  when '*'
    return a * c
  when '/'
    return a / c
  when '+'
    return a + c
  when '-'
    return a - c
  else
    print 'Unexpected operator: ' + b
    return ERROR
  end
end

def parse(s)
  r = []
  state = 0
  while s.length > 0
    case state
    when 0
      if s[0] == '('
        r << '('
        s = s[1..-1]
      else
        n = s[/^[-]?[0-9]+([.][0-9]+)?/]
        if s[0] == '-' && s[1] == '('
          r << '['
          s = s[2..-1]
        elsif !n
          print 'Input error: expected number, opening bracket, or unary -, got ' + s
          return ERROR
        else
          r << n.to_f
          s = s[n.length..-1]
          state = 1
        end
      end
    when 1
      if '^*/+-'.include? s[0]
        r << s[0]
        s = s[1..-1]
        state = 0
        if r.length > 3 && r[-2] != ')' && r[-3] != '(' && r[-3] != '[' && r[-4] != ')'
          unless '^*/+-'.include? r[-3]
            print 'Expected operator, got ' + r[-3]
            return ERROR
          end
          if P[r[-3]] >= P[r[-1]]
            r[-4] = operate(r[-4], r[-3], r[-2])
            return ERROR if r[-4] == ERROR
            r = r[0..-4] << r[-1]
          end
        end
      elsif s[0] == ')'
        r << ')'
        s = s[1..-1]
        while r.length > 3 && r[-2] != ')' && r[-3] != '(' && r[-3] != '[' && r[-4] != ')'
          r[-4] = operate(r[-4], r[-3], r[-2])
          return ERROR if r[-4] == ERROR
          r = r[0..-4] << r[-1]
        end
        if r.length > 3 && r[-3] == '('
          r = r[0...-3] << r[-2]
        elsif r.length > 3 && r[-3] == '['
          r = r[0...-3] << -r[-2]
        elsif r.length == 3 && r[0] == '('
          r = [r[1]]
        elsif r.length == 3 && r[0] == '['
          r = [-r[1]]
        end
      else
        print 'Input error: expected operator or closing bracket, got ' + s
        return ERROR
      end
    end
  end
  while r.length > 1
    r[-3] = operate(r[-3], r[-2], r[-1])
    return ERROR if r[-3] == ERROR
    r = r[0..-3]
  end
  r
end

File.open(ARGV[0]).each_line do |line|
  s = line.gsub(/\s+/, '')
  p = parse(s)
  if p == ERROR
    puts ' on input ' + s
  elsif p[0] && p[0].class == Float
    t = format '%.5f', p[0]
    puts t.sub(/[.]?0+$/, '')
  end
end
