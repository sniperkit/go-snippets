import System.Environment (getArgs)

main :: IO ()
main = do
    let stairs = 1 : 1 : zipWith (+) stairs (tail stairs)
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . (stairs!!) . read) $ [x | x <- lines input, not (null x)]
