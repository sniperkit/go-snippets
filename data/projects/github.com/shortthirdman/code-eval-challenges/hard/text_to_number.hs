import System.Environment (getArgs)
import Data.List (elemIndex)
import Data.Maybe (fromJust, isJust)

d0 :: [String]
d0 = ["zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"]
d1 :: [String]
d1 = ["ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
      "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"]
v1 :: [Int]
v1 = [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 30, 40, 50, 60, 70, 80, 90]

ttn             :: Int -> Int -> Int -> [String] -> Int
ttn r h q []     = r + q + h
ttn r h q (x:xs) | x == "zero"     = ttn r h 0 xs
                 | isJust d0i      = ttn r (fromJust d0i) q xs
                 | isJust d1i      = ttn r h (v1 !! fromJust d1i + q) xs
                 | x == "hundred"  = ttn r 0 (h*100) xs
                 | x == "thousand" = ttn (r+(h+q)*1000) 0 0 xs
                 | x == "million"  = ttn ((h+q)*1000000) 0 0 xs
                 where d0i = elemIndex x d0
                       d1i = elemIndex x d1

texttonumber       :: [String] -> Int
texttonumber []     = 0
texttonumber (x:xs) | x == "negative" = - texttonumber xs
                    | otherwise       = ttn 0 0 0 (x:xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . texttonumber . words) $ lines input
