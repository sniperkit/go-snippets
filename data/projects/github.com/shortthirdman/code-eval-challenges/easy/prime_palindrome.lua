function div(a, b) return math.floor(a/b) end

primes = {2, 3, 5, 7, 11, 13}
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

function isPalindrome(a)
  if a % 10 == 0 then
    return a == 0
  end
  local r = 0
  while a > r do
    r = 10*r + a%10
    if a == r then return true end
    a = div(a, 10)
  end
  return r == a
end

local maxprime, nextPrime = 0, primeSeq()
local j = nextPrime()
while j < 1000 do
  if isPalindrome(j) then maxprime = j end
  j = nextPrime()
end
print(maxprime)
