import System.Environment (getArgs)
import Data.Bits ((.&.))
import Data.Char (digitToInt)
import Data.List.Split (splitOn)

digits :: [Int]
digits = [252, 96, 218, 242, 102, 182, 190, 224, 254, 246]

blc :: [Int] -> [Int] -> Bool
blc _ []          = True
blc (x:xs) (y:ys) | (.&.) x y == y = blc xs ys
                  | otherwise      = False

blcd :: [Int] -> [Int] -> Bool
blcd xs ys | blc xs ys             = True
           | length xs > length ys = blcd (tail xs) ys
           | otherwise             = False

b2i         :: Int -> String -> Int
b2i x ""     = x
b2i x (y:ys) = b2i (x * 2 + digitToInt y) ys

d2i :: [Int] -> String -> [Int]
d2i xs []     = xs
d2i xs (y:ys) | y == '.'  = d2i (init xs ++ [succ $ last xs]) ys
              | otherwise = d2i (xs ++ [digits !! digitToInt y]) ys

brokenlcd         :: [String] -> String
brokenlcd [xs, ys] | blcd (map (b2i 0) $ words xs) (d2i [] ys) = "1"
                   | otherwise                               = "0"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (brokenlcd . splitOn ";") $ lines input
