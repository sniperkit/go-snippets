import System.Environment (getArgs)

maxlen      :: String -> String -> String
maxlen xs ys | length xs < length ys = ys
             | otherwise             = xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (foldl maxlen "" . words) $ lines input
