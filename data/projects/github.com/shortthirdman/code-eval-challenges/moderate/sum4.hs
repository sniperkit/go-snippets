import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (sort)

numzero      :: Int -> Int -> [Int] -> Int
numzero 0 0 _ = 1
numzero 0 _ _ = 0
numzero c z s | length s < c       = 0
              | z + c * head s > 0 = 0
              | z + c * last s < 0 = 0
              | otherwise          = numzero (c - 1) (z + head s) (tail s) + numzero c z (tail s)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . numzero 4 0 . sort . map read . splitOn ",") $ lines input
