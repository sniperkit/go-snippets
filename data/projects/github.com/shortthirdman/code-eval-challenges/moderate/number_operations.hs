import System.Environment (getArgs)
import Data.Char (chr, ord)
import Data.List (permutations, sort)

nmb             :: Int -> Int -> Int -> String -> Bool
nmb x _ _ []     = x == 42
nmb x y z (a:as) | y == 0    = nmb (x + ord a) (mod z 3) (div z 3) as
                 | y == 1    = nmb (x - ord a) (mod z 3) (div z 3) as
                 | otherwise = nmb (x * ord a) (mod z 3) (div z 3) as

numb     :: Int -> String -> Bool
numb 81 _ = False
numb x ys | nmb (ord (head ys)) (mod x 3) (div x 3) (tail ys) = True
          | otherwise                                         = numb (x+1) ys

numbo         :: [String] -> Bool
numbo []       = False
numbo (xs:xss) | numb 0 xs = True
               | otherwise = numbo xss

numberoperations   :: [Int] -> String
numberoperations xs | numbo ys  = "YES"
                    | otherwise = "NO"
                    where ys = permutations . map chr $ sort xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (numberoperations . map read . words) $ lines input
