import System.Environment (getArgs)

choose    :: (Integral a) => a -> a -> a
choose n k = foldl (\x i -> div (x * (n-i+1)) i) 1 [1..k]

luckytickets    :: Integer -> Integer -> Integer
luckytickets x y | 2 * x == y   = 0
                 | mod x 2 == 0 = b + luckytickets (succ x) y
                 | otherwise    = (-b) + luckytickets (succ x) y
                 where b = choose y x * choose (div (11 * y) 2 - 1 - 10 * x) (y - 1)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . luckytickets 0 . read) $ lines input
