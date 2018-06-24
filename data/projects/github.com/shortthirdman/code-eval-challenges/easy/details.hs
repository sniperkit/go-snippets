import System.Environment (getArgs)
import Data.List (isInfixOf)
import Data.List.Split (splitOn)

count :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

details       :: [String] -> Int
details []     = 10
details (x:xs) | isInfixOf "XY" x = 0
               | otherwise        = min (count '.' x) (details xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . details . splitOn ",") $ lines input
