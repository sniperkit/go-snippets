import System.Environment (getArgs)

swapnum   :: String -> String
swapnum xs = last xs : init (tail xs) ++ [head xs]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map swapnum . words) $ lines input
