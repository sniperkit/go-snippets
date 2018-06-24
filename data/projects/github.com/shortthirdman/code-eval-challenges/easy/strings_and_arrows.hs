import System.Environment (getArgs)

arrows   :: String -> Int
arrows xs | length xs < 5        = 0
          | take 5 xs == ">>-->" = 1 + arrows (drop 4 xs)
          | take 5 xs == "<--<<" = 1 + arrows (drop 4 xs)
          | otherwise            = arrows (tail xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . arrows) $ lines input
