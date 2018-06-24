import System.Environment (getArgs)

floors    :: Int -> Int -> Int
floors _ 0 = 0
floors _ 1 = 1
floors 0 _ = 0
floors 1 s = s
floors e s | e > s     = floors s s
           | otherwise = floors (e-1) (s-1) + floors e (s-1) + 1

tryFloor     :: Int -> [Int] -> Int
tryFloor n xs | floors e n < s = tryFloor (succ n) xs
              | otherwise      = n
              where e = head xs
                    s = last xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . tryFloor 0 . map read . words) $ lines input
