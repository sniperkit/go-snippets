import System.Environment (getArgs)
import Data.List.Split (splitOn)

within         :: [Int] -> Int -> Int -> Int -> Bool
within xs r p q = xs!!p >= xs!!r && xs!!p <= xs!!(r+2) && xs!!q >= xs!!(r+3) && xs!!q <= xs!!(r+1)

overlapping   :: [Int] -> String
overlapping xs | within xs 0 4 5 || within xs 0 6 7 || within xs 0 4 7 || within xs 0 6 5 || within xs 4 0 1 || within xs 4 2 3 || within xs 4 0 3 || within xs 4 2 1 = "True"
               | otherwise = "False"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (overlapping . map read . splitOn ",") $ lines input
