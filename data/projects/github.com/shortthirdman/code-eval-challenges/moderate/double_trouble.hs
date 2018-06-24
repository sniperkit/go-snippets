import System.Environment (getArgs)

doublet             :: (String, String) -> Int
doublet ("",   "")   = 1
doublet (x:xs, y:ys) | x == 'A' && y == 'B' = 0
                     | x == 'B' && y == 'A' = 0
                     | x == '*' && y == '*' = 2 * doublet (xs, ys)
                     | otherwise            = doublet (xs, ys)

doubletrouble   :: String -> String
doubletrouble xs | odd x     = "0"
                 | otherwise = show (doublet (splitAt (div x 2) xs))
                 where x = length xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map doubletrouble $ lines input
