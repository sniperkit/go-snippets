def robot(f, x, y)
  return 1 if x == 3 && y == 3
  ret = 0
  if x > 0 && f & (1 << (x - 1 + 4 * y)) == 0
    ret += robot(f | (1 << (x - 1 + 4 * y)), x - 1, y)
  end
  if y > 0 && f & (1 << (x + 4 * (y - 1))) == 0
    ret += robot(f | (1 << (x + 4 * (y - 1))), x, y - 1)
  end
  if x < 3 && f & (1 << (x + 1 + 4 * y)) == 0
    ret += robot(f | (1 << (x + 1 + 4 * y)), x + 1, y)
  end
  if y < 3 && f & (1 << (x + 4 * (y + 1))) == 0
    ret += robot(f | (1 << (x + 4 * (y + 1))), x, y + 1)
  end
  ret
end

puts robot(1, 0, 0)
