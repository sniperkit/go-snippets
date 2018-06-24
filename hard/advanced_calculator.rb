P = { '^' => 4, '%' => 3, '*' => 2, '/' => 2, '+' => 1, '-' => 1 }.freeze
ERROR = -999_999
O = ['(', '-(', '|', '-|', 'sqrt(', 'cos(', 'sin(', 'tan(', 'lg(', 'ln('].freeze

def factorial(a)
  if a < 0
    print 'Negative factorial: ' + a.to_s
    return ERROR
  end
  r = 1
  while a > 1
    r *= a
    a -= 1
  end
  r
end

def operate(a, b, c)
  case b
  when '^' then return a**c
  when '%' then return a % c
  when '*' then return a * c
  when '/' then return a / c
  when '+' then return a + c
  when '-' then return a - c
  end
  print 'Unexpected operator: ' + b
  ERROR
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
          r << '-('
          s = s[2..-1]
        elsif s[0] == '-' && s[1] == '|'
          r << '-|'
          s = s[2..-1]
        elsif s[0] == 'e'
          r << Math::E
          s = s[1..-1]
          state = 1
        elsif s[0] == 'P' && s[1] == 'i'
          r << Math::PI
          s = s[2..-1]
          state = 1
        elsif s[0] == '|'
          r << '|'
          s = s[1..-1]
        elsif s[0..4] == 'sqrt('
          r << 'sqrt('
          s = s[5..-1]
        elsif s[0..3] == 'cos('
          r << 'cos('
          s = s[4..-1]
        elsif s[0..3] == 'sin('
          r << 'sin('
          s = s[4..-1]
        elsif s[0..3] == 'tan('
          r << 'tan('
          s = s[4..-1]
        elsif s[0..2] == 'lg('
          r << 'lg('
          s = s[3..-1]
        elsif s[0..2] == 'ln('
          r << 'ln('
          s = s[3..-1]
        elsif !n
          print 'Input error: expected number, opening bracket,'\
                'or unary -, got ' + s
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
        if r.length > 3 && r[-2] != ')' && !(O.include? r[-3]) && r[-4] != ')'
          unless '^%*/+-'.include? r[-3]
            print 'Expected operator, got ' + r[-3]
            return ERROR
          end
          if P[r[-3]] >= P[r[-1]]
            r[-4] = operate(r[-4], r[-3], r[-2])
            return ERROR if r[-4] == ERROR
            r = r[0..-4] << r[-1]
          end
        end
      elsif s[0..2] == 'mod'
        r << '%'
        s = s[3..-1]
        state = 0
        if r.length > 3 && r[-2] != ')' && !(O.include? r[-3]) && r[-4] != ')'
          unless '^%*/+-'.include? r[-3]
            print 'Expected operator, got ' + r[-3]
            return ERROR
          end
          if P[r[-3]] >= P[r[-1]]
            r[-4] = operate(r[-4], r[-3], r[-2])
            return ERROR if r[-4] == ERROR
            r = r[0..-4] << r[-1]
          end
        end
      elsif s[0] == '!'
        r[-1] = factorial(r[-1])
        return ERROR if r[-1] == ERROR
        s = s[1..-1]
      elsif s[0] == ')'
        r << ')'
        s = s[1..-1]
        while r.length > 3 && r[-2] != ')' &&
              !(O.include? r[-3]) && r[-4] != ')'
          r[-4] = operate(r[-4], r[-3], r[-2])
          return ERROR if r[-4] == ERROR
          r = r[0..-4] << r[-1]
        end
        if r.length > 3 && r[-3] == '('
          if s[0] == '!'
            s = s[1..-1]
            r[-2] = factorial(r[2])
            return ERROR if r[-1] == ERROR
          end
          r = r[0...-3] << r[-2]
        elsif r.length > 3 && r[-3] == '-('
          if s[0] == '!'
            s = s[1..-1]
            r[-2] = factorial(r[2])
            return ERROR if r[-1] == ERROR
          end
          r = r[0...-3] << -r[-2]
        elsif r.length > 3 && r[-3] == 'sqrt('
          r = r[0...-3] << Math.sqrt(r[-2])
        elsif r.length > 3 && r[-3] == 'sin('
          r = r[0...-3] << Math.sin(r[-2])
        elsif r.length > 3 && r[-3] == 'cos('
          r = r[0...-3] << Math.cos(r[-2])
        elsif r.length > 3 && r[-3] == 'tan('
          r = r[0...-3] << Math.tan(r[-2])
        elsif r.length > 3 && r[-3] == 'ln('
          r = r[0...-3] << Math.log(r[-2])
        elsif r.length > 3 && r[-3] == 'lg('
          r = r[0...-3] << Math.log10(r[-2])
        elsif r.length == 3 && r[0] == '('
          if s[0] == '!'
            s = s[1..-1]
            r[1] = factorial(r[1])
            return ERROR if r[-1] == ERROR
          end
          r = [r[1]]
        elsif r.length == 3 && r[0] == '-('
          if s[0] == '!'
            s = s[1..-1]
            r[1] = factorial(r[1])
            return ERROR if r[-1] == ERROR
          end
          r = [-r[1]]
        elsif r.length == 3 && r[0] == 'sqrt('
          r = [Math.sqrt(r[1])]
        elsif r.length == 3 && r[0] == 'sin('
          r = [Math.sin(r[1])]
        elsif r.length == 3 && r[0] == 'cos('
          r = [Math.cos(r[1])]
        elsif r.length == 3 && r[0] == 'tan('
          r = [Math.tan(r[1])]
        elsif r.length == 3 && r[0] == 'ln('
          r = [Math.log(r[1])]
        elsif r.length == 3 && r[0] == 'lg('
          r = [Math.log10(r[1])]
        end
      elsif s[0] == '|'
        r << '|'
        s = s[1..-1]
        while r.length > 3 && r[-2] != ')' &&
              !(['|', '-|'].include? r[-3]) && r[-4] != ')'
          r[-4] = operate(r[-4], r[-3], r[-2])
          return ERROR if r[-4] == ERROR
          r = r[0..-4] << r[-1]
        end
        if r.length > 3 && r[-3] == '|'
          if s[0] == '!'
            s = s[1..-1]
            r[-2] = factorial(r[2])
            return ERROR if r[-1] == ERROR
          end
          r = r[0...-3] << r[-2].abs
        elsif r.length > 3 && r[-3] == '-|'
          if s[0] == '!'
            s = s[1..-1]
            r[-2] = factorial(r[2])
            return ERROR if r[-1] == ERROR
          end
          r = r[0...-3] << -r[-2].abs
        elsif r.length == 3 && r[0] == '|'
          r = [r[1].abs]
        elsif r.length == 3 && r[0] == '-|'
          if s[0] == '!'
            s = s[1..-1]
            r[1] = factorial(r[1].abs)
            return ERROR if r[-1] == ERROR
          end
          r = [-r[1].abs]
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
  if p == ERROR then puts ' on input ' + s
  elsif p[0] && p[0].class == Float
    t = format '%.5f', p[0]
    puts t.sub(/[.]?0+$/, '')
  end
end
