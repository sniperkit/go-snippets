import System.Environment (getArgs)
import Data.List.Split (splitOn)

sLines :: Int -> Int -> Int -> [[Int]] -> Int
sLines i j k xss | i >= length xss - 2               = 0
                 | j >= length xss - 1               = sLines (i + 1) (i + 2) 0 xss
                 | k >= length xss                   = sLines i (j + 1) 0 xss
                 | k == i || k == j                  = sLines i j (k + 1) xss
                 | (dx * dy' == dx' * dy) && (k > j) = 1 + sLines i (j + 1) 0 xss
                 | dx * dy' == dx' * dy              = sLines i (j + 1) 0 xss
                 | otherwise                         = sLines i j (k + 1) xss
                 where dx  = head (xss!!i) - head (xss!!j)
                       dy  = last (xss!!i) - last (xss!!j)
                       dx' = head (xss!!i) - head (xss!!k)
                       dy' = last (xss!!i) - last (xss!!k)

straightLines :: [[Int]] -> Int
straightLines = sLines 0 1 2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . straightLines . map (map read . words) . splitOn " | ") $ lines input
