import System.Environment (getArgs)
import Data.List.Split (splitOn)

largest           :: Int -> Int -> [Int] -> String
largest i _ []     = show i
largest i j (x:xs) = largest k l xs
               where k = max (x + j) i
                     l = max (x + j) 0

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (largest (minBound::Int) 0 . map read . splitOn ",") $ lines input
