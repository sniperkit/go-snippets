import System.Environment (getArgs)

fizz        :: Int -> [String] -> [Int] -> [String]
fizz i xs ys | i > last ys                  = reverse xs
             | mod i f == 0 && mod i b == 0 = fizz (succ i) ("FB" : xs) ys
             | mod i f == 0                 = fizz (succ i) ("F" : xs) ys
             | mod i b == 0                 = fizz (succ i) ("B" : xs) ys
             | otherwise                    = fizz (succ i) (show i : xs) ys
             where f = head ys
                   b = ys!!1

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . fizz 1 [] . map read . words) $ lines input
