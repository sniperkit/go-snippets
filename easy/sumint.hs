import System.Environment (getArgs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    print . sum . map read $ lines input
