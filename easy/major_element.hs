import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (nub)

count     :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

majE          :: [Int] -> [Int] -> String
majE []     _  = "None"
majE (x:xs) ys | count x ys > div (length ys) 2 = show x
               | otherwise                      = majE xs ys

majorElement  :: [Int] -> String
majorElement s = majE (nub s) s

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (majorElement . map read . splitOn ",") $ lines input
