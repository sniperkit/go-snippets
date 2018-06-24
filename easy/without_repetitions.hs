import System.Environment (getArgs)

wr          :: String -> String -> String
wr xs []     = xs
wr xs [y]    = xs ++ [y]
wr xs (y:ys) | y == head ys = wr xs ys
             | otherwise    = wr (xs ++ [y]) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (wr "") $ lines input
