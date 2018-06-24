import System.Environment (getArgs)
import Data.List (transpose)

splitEvery     :: Int -> [a] -> [[a]]
splitEvery _ [] = []
splitEvery x ys = take x ys : splitEvery x (drop x ys)

sq    :: Int -> Int -> Int
sq x y | x*x >= y  = x
       | otherwise = sq (x+1) y
sqr :: Int -> Int
sqr  = sq 1

rotate   :: [String] -> [String]
rotate xs = concatMap reverse . transpose $ splitEvery (sqr $ length xs) xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . rotate . words) $ lines input
