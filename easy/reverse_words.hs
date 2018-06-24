import System.Environment (getArgs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . reverse . words) $ [x | x <- lines input, not (null x)]
