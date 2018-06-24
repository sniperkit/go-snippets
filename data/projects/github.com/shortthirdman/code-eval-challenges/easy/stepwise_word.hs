import System.Environment (getArgs)

explode   :: String -> [String]
explode [] = []
explode xs = explode (init xs) ++ [replicate (pred (length xs)) '*' ++ [last xs]]

maxlen      :: String -> String -> String
maxlen xs ys | length xs < length ys = ys
             | otherwise             = xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . explode . foldl maxlen "" . words) $ lines input
