import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (transpose)

findHighest    :: [[Int]] -> [Int]
findHighest xss = map maximum $ transpose xss

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . findHighest . map (map read . words) . splitOn "|") $ lines input
