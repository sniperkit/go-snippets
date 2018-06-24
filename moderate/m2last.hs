import System.Environment (getArgs)

nonempty  :: [[a]] -> [[a]]
nonempty s = [x | x <- s, not (null x)]

m2last  :: [String] -> String
m2last s | l > length s = ""
         | otherwise    = s !! (length s - l)
         where l = succ . read $ last s

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . nonempty . map (m2last . words) $ lines input
