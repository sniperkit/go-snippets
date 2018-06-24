import System.Environment (getArgs)
import Data.Bits ((.&.), popCount, shiftL)
import Data.Char (digitToInt)
import Data.List.Split (splitOn)

numCar        :: Int -> Int -> [Int] -> Int
numCar rm i zs | i == rm   = 0
               | otherwise = popCount ((.&.) (zs!!i) (pred (shiftL 1 rm))) + numCar rm (succ i) zs

isNum                  :: Int -> Int -> Int -> Int -> Int -> Int -> [Int] -> Bool
isNum rn rm x y k rk zs | k == rm   = rn == rk
                        | rk > rn   = False
                        | otherwise = isNum rn rm x y (succ k) (rk + popCount mask) zs
                        where mask = (.&.) (zs!!(x+k)) (shiftL (pred (shiftL 1 rm)) y)

isbw             :: Int -> Int -> Int -> Int -> [Int] -> (Int, Int)
isbw rn rm x y zs | z == rm                = (rm, sum (map popCount zs))
                  | x == z - rm + 1        = (rm, rn)
                  | (x == 0) && (y == 0)   = isbw (numCar rm 0 zs) rm 0 1 zs
                  | y == z - rm + 1        = isbw rn rm (succ x) 0 zs
                  | isNum rn rm x y 0 0 zs = isbw rn rm x (succ y) zs
                  | otherwise              = isbw 0 (succ rm) 0 0 zs
                  where z = length zs

b2i         :: Int -> String -> Int
b2i x ""     = x
b2i x (y:ys) = b2i (x * 2 + digitToInt y) ys

blackOrWhite :: [String] -> String
blackOrWhite xs = show (fst b) ++ "x" ++ show (fst b) ++ ", " ++ show (snd b)
                where b = isbw 0 1 0 0 (map (b2i 0) xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (blackOrWhite . splitOn " | ") $ lines input
