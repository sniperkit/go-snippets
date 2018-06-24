import System.Environment (getArgs)
import Data.List (sortBy)

jolly   :: [Int] -> String
jolly [] = "Jolly"
jolly xs | head xs /= length xs = "Not jolly"
         | otherwise            = jolly (tail xs)

delta          :: [Int] -> [Int] -> [Int]
delta _  []     = []
delta xs [_]    = xs
delta xs (y:ys) = delta (abs (y - head ys) : xs) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (jolly . sortBy (flip compare) . delta [] . map read . tail . words) $ lines input
