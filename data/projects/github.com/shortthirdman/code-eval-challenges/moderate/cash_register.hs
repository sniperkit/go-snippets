import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate)

units :: [String]
units  = ["PENNY", "NICKEL", "DIME", "QUARTER", "HALF DOLLAR", "ONE",
          "TWO", "FIVE", "TEN", "TWENTY", "FIFTY", "ONE HUNDRED"]
value :: [Int]
value  = [1, 5, 10, 25, 50, 100,
          200, 500, 1000, 2000, 5000, 10000]

centi :: String -> Int
centi xs | elem '.' xs = 100 * read (head xs') + read (last xs')
         | otherwise   = 100 * read xs
         where xs' = splitOn "." xs

cash              :: [Int] -> [Int] -> [String]
cash []     []     = []
cash (x:xs) (y:ys) | y == 0 = cash xs ys
                   | otherwise = units!!x : cash (x:xs) ((y-1):ys)

cashr       :: [Int] -> Int -> Int -> String
cashr xs _ 0 = intercalate "," (reverse (cash [0..11] xs))
cashr xs y z | z >= value!!y = cashr (take y xs ++ [xs!!y + 1] ++ drop (y+1) xs) y (z - value!!y)
             | otherwise     = cashr xs (y-1) z

cashreg         :: [String] -> String
cashreg [xs, ys] | x > y     = "ERROR"
                 | x == y    = "ZERO"
                 | otherwise = cashr (replicate 12 0) 11 (y - x)
                 where x = centi xs
                       y = centi ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (cashreg . splitOn ";") $ lines input
