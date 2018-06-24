import System.Environment (getArgs)

dsquare           :: Int -> Int -> [Int] -> Int -> Int -> Int
dsquare j n xs t x | j >= t = n
                   | dsqr u && notElem u xs = dsquare (j+1) (n+1) (j*j : xs) t x
                   | otherwise = dsquare (j+1) n xs t x
                   where u = x - j * j

dsqr  :: Int -> Bool
dsqr x = s * s == x
       where s = sqr x

sq    :: Int -> Int -> Int
sq x y | x * x > y = x - 1
       | otherwise = sq (x + 1) y

sqr  :: Int -> Int
sqr = sq 0

dsquares  :: Int -> Int
dsquares x = dsquare (sqr (div x 2)) 0 [] (sqr x + 1) x

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . dsquares . read) . tail $ lines input
