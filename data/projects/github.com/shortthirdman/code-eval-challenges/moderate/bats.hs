import System.Environment (getArgs)
import Data.List (sort)

bat                  :: Int -> Int -> Int -> Int -> Int -> Int -> Int -> [Int] -> Int
bat l d n c s t tx xs | s > l - 6            = c
                      | s > tx - d && t == n = bat l d n c (tx+d) t (l-6+d) xs
                      | s > tx - d           = bat l d n c (tx+d) (t+1) (head xs) (tail xs)
                      | otherwise            = bat l d n (c+1) (s+d) t tx xs

bats           :: [Int] -> Int
bats (a:d:n:xs) = bat a d n 0 6 0 (6-d) (sort xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . bats . map read . words) $ lines input
