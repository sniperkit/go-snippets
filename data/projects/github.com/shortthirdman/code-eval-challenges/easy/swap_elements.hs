import System.Environment (getArgs)

swapel       :: Int -> Int -> [String] -> [String]
swapel a b xs = take a xs ++ [xs!!b] ++ take (b-a-1) (drop (a+1) xs) ++ [xs!!a] ++ drop (b+1) xs

swe          :: [String] -> [String] -> [String]
swe xs []     = xs
swe xs (y:ys) | v == w    = swe xs ys
              | otherwise = swe (swapel a b xs) ys
              where v = read (takeWhile (/= '-') y)
                    w = read (drop 1 (takeWhile (/= ',') (dropWhile (/= '-') y)))
                    a = min v w
                    b = max v w

swapelements   :: [String] -> [String]
swapelements xs = swe ys zs
                where ys = takeWhile (/= ":") xs
                      zs = drop 1 (dropWhile (/= ":") xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . swapelements . words) $ lines input
