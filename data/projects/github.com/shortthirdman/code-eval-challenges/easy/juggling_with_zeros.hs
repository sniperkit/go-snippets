import System.Environment (getArgs)

juggle         :: Integer -> [Int] -> Integer
juggle i []     = i
juggle i (x:xs) | x == 1    = juggle (i * 2 ^ head xs) (tail xs)
                | otherwise = juggle ((succ i * 2 ^ head xs) - 1) (tail xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . juggle 0 . map length . words) $ lines input
