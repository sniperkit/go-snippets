import System.Environment (getArgs)
import Data.Bits ((.&.), (.|.), bit)
import Data.Char (digitToInt)
import Data.List (sort)

testmagic         :: Int -> Int -> Int -> String -> Bool
testmagic x y z xs | x == length xs && y == 0 = True
                   | x == length xs           = False
                   | (.&.) z (bit p) > 0      = False
                   | otherwise                = testmagic (succ x) p ((.|.) z (bit p)) xs
                   where p = mod (y + digitToInt (xs!!y)) (length xs)

notuni :: String -> Bool
notuni [_]    = False
notuni (x:xs) | x == head xs = True
              | otherwise    = notuni xs

notuniq  :: Int -> Bool
notuniq x = notuni . sort $ show x

contain0  :: Int -> Bool
contain0 x | x < 10        = False
           | mod x 10 == 0 = True
           | otherwise     = contain0 (div x 10)

ismagic :: Int -> Bool
ismagic x | contain0 x = False
          | notuniq  x = False
          | otherwise  = testmagic 0 0 0 (show x)

magnumb    :: Int -> Int -> [Int]
magnumb x y | x > y     = []
            | ismagic x = x : magnumb (succ x) y
            | otherwise = magnumb (succ x) y

magicnumbers       :: [Int] -> [Int]
magicnumbers [x, y] | null (magnumb x y) = [-1]
                    | otherwise          = magnumb x y

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . magicnumbers . map read . words) $ lines input
