import System.Environment (getArgs)
import Data.Char (isDigit, isSpace)
import Data.List (intercalate, sort)

roadTrip            :: Int -> [Int] -> [Int] -> [Int]
roadTrip _ xs []     = xs
roadTrip _ [] (y:ys) = roadTrip y [y] ys
roadTrip w xs (y:ys) = roadTrip y (xs ++ [y - w]) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . map show . roadTrip 0 []. sort . map read . words) . lines $ [x | x <- input, isDigit x || isSpace x]
