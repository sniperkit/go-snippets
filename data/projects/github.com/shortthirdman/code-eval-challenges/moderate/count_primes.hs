import System.Environment (getArgs)
import Data.List.Split (splitOn)

factors  :: Int -> [Int]
factors n = [x | x <- [1..n], mod n x == 0]

prime  :: Int -> Bool
prime x = factors x == [1, x]

countPrimes       :: [Int] -> Int
countPrimes []     = -1
countPrimes [a, b] = length [x | x <- [a..b], prime x]
countPrimes _      = -1

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . countPrimes . map read . splitOn ",") $ lines input
