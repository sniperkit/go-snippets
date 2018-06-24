import System.Environment (getArgs)
import Data.Char (digitToInt)
import Data.List (elemIndex)
import Data.List.Split (splitOn)
import Data.Maybe (fromJust)

alpha :: String
alpha = " !\"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

gron              :: String -> String -> String
gron _      []     = []
gron (x:xs) (y:ys) = alpha !! n : gron xs ys
                   where m = fromJust (elemIndex y alpha) - digitToInt x
                         n | m >= 0    = m
                           | otherwise = m + length alpha

gronsfeld         :: [String] -> String
gronsfeld [xs, ys] = gron (cycle xs) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (gronsfeld . splitOn ";") $ lines input
