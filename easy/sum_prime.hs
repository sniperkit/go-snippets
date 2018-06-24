factors  :: Int -> [Int]
factors n = [x | x <- [1..n], mod n x == 0]

prime  :: Int -> Bool
prime x = factors x == [1, x]

main :: IO ()
main = print . sum $ take 1000 [x | x <- 2 : [3, 5..], prime x]
