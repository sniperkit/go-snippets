factors  :: Int -> [Int]
factors n = [x | x <- [1..n], mod n x == 0]

prime  :: Int -> Bool
prime x = factors x == [1, x]

pali  :: Int -> Bool
pali x | x < 10                             = True
       | x < 100 && (div x 10 == mod x 10)  = True
       | x > 100 && (div x 100 == mod x 10) = True
       | otherwise                          = False

main :: IO ()
main = print $ last [x | x <- 2 : [3, 5..1000], pali x, prime x]
