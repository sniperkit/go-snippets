import System.Environment (getArgs)

pasc      :: Int -> Int -> Int -> String
pasc r j i | j == i    = ""
           | otherwise = " " ++ show r' ++ pasc r' (succ j) i
           where r' = div (r*(i-j)) (succ j)

pascal    :: Int -> Int -> String
pascal 0 y = "1" ++ pascal 1 y
pascal i y | i == y    = ""
           | otherwise = " 1" ++ pasc 1 0 i ++ pascal (succ i) y

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (pascal 0 . read) $ lines input
