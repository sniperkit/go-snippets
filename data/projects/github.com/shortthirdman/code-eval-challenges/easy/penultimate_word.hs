import System.Environment (getArgs)

penulti  :: [String] -> String
penulti s | length s < 2 = ""
          | otherwise    = s !! (length s - 2)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (penulti . words) $ lines input
