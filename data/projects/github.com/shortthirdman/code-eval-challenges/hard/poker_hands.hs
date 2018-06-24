import System.Environment (getArgs)
import Data.List (sort)
import Data.Bits (shiftL)

numCards :: Int
numCards = 5

cardval :: Char -> Int
cardval '2' = 2
cardval '3' = 3
cardval '4' = 4
cardval '5' = 5
cardval '6' = 6
cardval '7' = 7
cardval '8' = 8
cardval '9' = 9
cardval 'T' = 10
cardval 'J' = 11
cardval 'Q' = 12
cardval 'K' = 13
cardval 'A' = 14
cardval _   = -1

wheel                 :: [Int] -> [Int]
wheel [2, 3, 4, 5, 14] = [1, 2, 3, 4, 5]
wheel xs               = xs

rank     :: Int -> [Int] -> Int
rank x [] = x
rank x ys = rank (shiftL x 4 + last ys) (init ys)

cranks   :: [String] -> [Int]
cranks xs = wheel (sort [cardval (head x) | x <- xs])

cardranks   :: [String] -> Int
cardranks xs = rank 0 (cranks xs)

straight  :: Int -> Bool
straight x = x0 + 1 == x1 && x1 + 1 == x2 && x2 + 1 == x3 && x3 + 1 == x4
           where x0 = mod x 16
                 x1 = mod (div x 16) 16
                 x2 = mod (div x 256) 16
                 x3 = mod (div x 4096) 16
                 x4 = div x 65536

flush         :: [String] -> Bool
flush [xs, ys] | last xs == last ys       = True
               | otherwise                = False
flush (x:xs)   | last x /= last (head xs) = False
               | otherwise                = flush xs

count     :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

kind     :: Int -> [Int] -> Int
kind x ys | x > length ys = 0
          | x == c        = y
          | otherwise     = kind x (drop c ys)
          where y = head ys
                c = count y ys

twopair  :: Int -> Int
twopair x | x0 == x1 && x2 == x3 && x0 /= x2 && x0 /= x4 && x2 /= x4 = x0
          | x0 == x1 && x3 == x4 && x0 /= x3 && x0 /= x2 && x2 /= x3 = x0
          | x1 == x2 && x3 == x4 && x1 /= x3 && x0 /= x1 && x0 /= x3 = x1
          | otherwise                                                = 0
          where x0 = mod x 16
                x1 = mod (div x 16) 16
                x2 = mod (div x 256) 16
                x3 = mod (div x 4096) 16
                x4 = div x 65536

handrank     :: [String] -> Int
handrank hand | straight ranks && flush hand         = shiftL 8 24 + ranks
              | kind 4 crnks > 0                     = shiftL 7 24 + shiftL (kind 4 crnks) 4  + kind 1 crnks
              | kind 3 crnks > 0 && kind 2 crnks > 0 = shiftL 6 24 + shiftL (kind 3 crnks) 4  + kind 2 crnks
              | flush hand                           = shiftL 5 24 + ranks
              | straight ranks                       = shiftL 4 24 + ranks
              | kind 3 crnks > 0                     = shiftL 3 24 + shiftL (kind 3 crnks) 20 + ranks
              | twopair ranks > 0                    = shiftL 2 24 + shiftL (twopair ranks) 4 + kind 1 crnks
              | kind 2 crnks > 0                     = shiftL 1 24 + shiftL (kind 2 crnks) 20 + ranks
              | otherwise                            = ranks
              where ranks = cardranks hand
                    crnks = cranks hand

poker  :: [String] -> String
poker xs | x > y     = "left"
         | x < y     = "right"
         | otherwise = "none"
         where x = handrank (take numCards xs)
               y = handrank (drop numCards xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (poker . words) $ lines input
