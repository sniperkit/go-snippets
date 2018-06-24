primes = {2, 3}
m = 5
function isPrime(a)
  local b = 1
  while a % primes[b] > 0 do
    if primes[b] * primes[b] > a then
      primes[#primes + 1] = a
      return true
    end
    b = b + 1
  end
  return false
end

function primeSeq()
  local p0, i = 1, 0
  return function()
           if p0 <= #primes then
             i, p0 = primes[p0], p0+1
             return i
           end
           i = i + 2
           while not isPrime(i) do
             i = i + 2
           end
           p0 = p0 + 1
           return i
         end
end

local nextPrime = primeSeq()
for line in io.lines(arg[1]) do
  local n = tonumber(line)
  while 2^m - 1 < n do
    isPrime(m)
    m = m + 2
  end
  local r = {}
  for i = 1, #primes do
    if 2^primes[i] - 1 >= n then break end
    r[#r + 1] = 2^primes[i] - 1
  end
  print(table.concat(r, ", "))
end
