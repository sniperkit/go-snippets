import System.Environment (getArgs)

spltn                   :: Int -> Int -> Int -> [String] -> Int
spltn a s b [[], []]     = a + b * s
spltn a s b [x:xs, y:ys] | y == '+'  = spltn (a + b * s)   1  0 [x:xs, ys]
                         | y == '-'  = spltn (a + b * s) (-1) 0 [x:xs, ys]
                         | otherwise = spltn a s (b*10 + read [x]) [xs, ys]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . spltn 0 1 0 . words) $ lines input
