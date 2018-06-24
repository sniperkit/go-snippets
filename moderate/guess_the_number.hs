import System.Environment (getArgs)

guess               :: Int -> Int -> [String] -> String
guess lo (-1) (x:xs) = guess lo (read x) xs
guess lo hi   (x:xs) | x == "Lower"  = guess lo (d - 1) xs
                     | x == "Higher" = guess (d + 1) hi xs
                     | otherwise     = show d
                     where c  = mod (lo+hi) 2
                           d  = div (lo+hi) 2 + c

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (guess 0 (-1) . words) $ lines input
