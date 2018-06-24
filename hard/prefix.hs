import System.Environment (getArgs)

prefix       :: [String] -> Float
prefix    [x] = (read :: String -> Float) x
prefix (x:xs) | x == "+"  = l + prefix (init xs)
              | x == "*"  = l * prefix (init xs)
              | x == "/"  = prefix (init xs) / l
              | otherwise = -1
              where l = (read :: String -> Float) (last xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . round . prefix . words) $ lines input
