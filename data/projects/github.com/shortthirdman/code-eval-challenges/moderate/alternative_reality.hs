import System.Environment (getArgs)

coins :: [Int]
coins = [50, 25, 10, 5, 1]

alternat         :: Int -> [Int] -> Int
alternat 0 _      = 1
alternat _ [1]    = 1
alternat x (y:ys) | y > x     = alternat x ys
                  | otherwise = alternat (x - y) (y:ys) + alternat x ys

alternative   :: String -> String
alternative xs = show $ alternat (read xs) coins

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map alternative $ lines input
