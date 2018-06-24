import System.Environment (getArgs)
import Data.Char (digitToInt, intToDigit)

count :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

selfd    :: Int -> String -> Bool
selfd i s | i > 9 || i == length s                      = True
          | digitToInt (s!!i) == count (intToDigit i) s = selfd (succ i) s
          | otherwise                                   = False

selfDesc  :: String -> String
selfDesc s | selfd 0 s = "1"
           | otherwise = "0"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map selfDesc $ lines input
