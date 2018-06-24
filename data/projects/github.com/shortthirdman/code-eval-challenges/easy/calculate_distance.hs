import System.Environment (getArgs)

stripChars :: String -> String -> String
stripChars = filter . flip notElem

distance                 :: [Double] -> String
distance [x1, y1, x2, y2] = show . floor $ sqrt (x*x + y*y)
                          where x = x1 - x2
                                y = y1 - y2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (distance . map read . words) . lines $ stripChars "()," input
