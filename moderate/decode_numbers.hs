import System.Environment (getArgs)

decode  :: String -> Int
decode []     = 1
decode [_]    = 1
decode (x:xs) | x /= '1' && x /= '2'                          = decode xs
              | head xs == '0' || (x == '2' && head xs > '6') = decode (tail xs)
              | otherwise                                     = decode xs + decode (tail xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . decode) $ lines input
