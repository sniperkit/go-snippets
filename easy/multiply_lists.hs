import System.Environment (getArgs)
import Data.List.Split (splitOn)

multiply   :: [[Int]] -> [Int]
multiply xs = [x * y | (x, y) <- zip (head xs) (last xs)]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . multiply . map (map read . words) . splitOn "|") $ lines input
