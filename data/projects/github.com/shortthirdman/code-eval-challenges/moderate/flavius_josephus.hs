import System.Environment (getArgs)
import Data.List.Split (splitOn)

flaviu     :: [Int] -> Int -> Int -> [Int]
flaviu [] _ _ = []
flaviu xs n m = cs!!(m-1) : flaviu (take (n-1) (drop m cs)) (n-1) m
            where cs = cycle xs

flavius       :: [Int] -> [Int]
flavius [n, m] = flaviu [0..(n-1)] n m

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . flavius . map read . splitOn ",") $ lines input
