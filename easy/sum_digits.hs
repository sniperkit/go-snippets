import System.Environment (getArgs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . sum . map (read . (:[]))) $ lines input
