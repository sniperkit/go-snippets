import System.Environment (getArgs)
import Data.List (intercalate)

prime     :: [Int] -> Int -> Bool
prime _ 2  = True
prime xs n | x*x > n      = True
           | mod n x == 0 = False
           | otherwise    = prime (tail xs) n
           where x = head xs

primesTo       :: [Int] -> Int -> Int -> [Int]
primesTo xs i n | i >= n           = xs
                | i == 2           = primesTo [2] 3 n
                | not (prime xs i) = primesTo xs (i + 2) n
                | otherwise        = primesTo (xs ++ [i]) (i + 2) n

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . map show . primesTo [] 2 . read) $ lines input
