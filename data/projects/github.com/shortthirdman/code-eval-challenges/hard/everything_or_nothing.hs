import System.Environment (getArgs)
import Data.Bits ((.&.), (.|.))
import Data.Char (digitToInt)

pm    :: Char -> Int
pm 'r' = 4
pm 'w' = 2
pm 'g' = 1

clean   :: String -> [Int]
clean xs | c == 'r'  = [a, b, 4]
         | c == 'w'  = [a, b, 2]
         | otherwise = [a, b, 1, d, e]
         where a = digitToInt (xs!!5) - 1
               b = digitToInt (xs!!13) - 1
               c = xs!!16
               d = pm (xs!!23)
               e = digitToInt (last xs) - 1

everyt            :: [Int] -> [[Int]] -> Bool
everyt xs []       = True
everyt xs (ys:yss) | (.&.) p a == 0 = False
                   | a == 1         = everyt (take j xs ++ [(.|.) (xs!!j) (ys!!3)] ++ drop (j + 1) xs) yss
                   | otherwise      = everyt xs yss
                   where a = ys!!2
                         p = xs!!(head ys * 3 + ys!!1)
                         j = last ys * 3 + ys!!1

everything    :: [String] -> String
everything xss | everyt perm (map clean xss) = "True"
               | otherwise                   = "False"
               where perm = [7, 3, 0, 6, 2, 4, 5, 1, 5, 3, 7, 1, 6, 0, 2, 4, 2, 6]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (everything . words) $ lines input
