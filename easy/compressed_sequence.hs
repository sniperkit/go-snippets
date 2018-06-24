import System.Environment (getArgs)

compress                :: [Int] -> [Int] -> [Int]
compress xs       []     = reverse xs
compress []       (y:ys) = compress [y, 1] ys
compress (x:n:xs) (y:ys) | x /= y    = compress (y:1:x:n:xs) ys
                         | otherwise = compress (x:(n+1):xs) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . compress [] . map read . words) $ lines input
