import System.Environment (getArgs)
import Data.List (sort)

mindist   :: [Int] -> Int
mindist xs = foldr (\ x -> (+) (abs (x - u))) 0 ys
           where ys = sort (tail xs)
                 u  = ys !! div (head xs) 2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . mindist . map read . words) $ lines input
