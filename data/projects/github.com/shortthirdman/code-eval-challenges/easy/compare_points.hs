import System.Environment (getArgs)

cardinal                 :: [Int] -> String
cardinal [x1, y1, x2, y2] | x1 == x2  = if y1 == y2 then "here" else if y1 < y2 then "N" else "S"
                          | y1 == y2  = if x1 < x2 then "E" else "W"
                          | x1 < x2   = if y1 < y2 then "NE" else "SE"
                          | otherwise = if y1 < y2 then "NW" else "SW"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (cardinal . map read . words) $ lines input
