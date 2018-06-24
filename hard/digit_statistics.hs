import System.Environment (getArgs)

stats :: [Int]
stats = [2, 4, 8, 6, 3, 9, 7, 1, 4, 6, 4, 6, 5, 5, 5, 5, 6, 6, 6, 6, 7, 9, 3, 1, 8, 4, 2, 6, 9, 1, 9, 1]

dig2         :: Int -> Int -> Integer -> [Integer] -> [Integer]
dig2 i a n xs | i == fromIntegral (mod n 4) = xs
              | otherwise                   = dig2 (i+1) a n (take j xs ++ [xs!!j + 1] ++ drop (j+1) xs)
              where j = stats!!(a*4+i)

dig1         :: Int -> Int -> Integer -> [Integer] -> [Integer]
dig1 4 _ _ xs = xs
dig1 i a m xs = dig1 (i+1) a m (take j xs ++ [xs!!j + m] ++ drop (j+1) xs)
              where j = stats!!(a*4+i)

beautify       :: [Integer] -> String
beautify [x]    = "9: " ++ show x
beautify (x:xs) = show (9 - length xs) ++ ": " ++ show x ++ ", " ++ beautify xs

digitStat       :: [Integer] -> String
digitStat [b, n] = beautify (dig2 0 a n (dig1 0 a m (replicate 10 0)))
                 where a = fromIntegral (b - 2)
                       m = div n 4

main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (digitStat . map read . words) $ lines input
