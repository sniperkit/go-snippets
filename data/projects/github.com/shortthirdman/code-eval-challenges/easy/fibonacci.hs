import System.Environment (getArgs)

fibs :: [Integer]
fibs = 0 : 1 : zipWith (+) fibs (tail fibs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . (fibs!!) . read) $ lines input
