import System.Environment (getArgs)
import Data.Char (chr, isAlpha, toLower)
import Data.List (sortBy)

count     :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

countChars       :: Int -> [Int] -> String -> [Int]
countChars i is s | i == 27   = is
                  | otherwise = countChars (i+1) (count (chr (97+i)) s : is) s

beauty       :: Int -> Int -> [Int] -> String
beauty i x xs | i == 0 || head xs == 0 = show x
              | otherwise              = beauty (i-1) (head xs * i + x) (tail xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (beauty 26 0 . sortBy (flip compare) . countChars 0 []) $ lines [toLower x | x <- input, isAlpha x || x == '\n']
