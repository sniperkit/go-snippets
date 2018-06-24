import System.Environment (getArgs)

parent   :: Int -> Int
parent 3  = 8
parent 8  = 30
parent 10 = 20
parent 20 = 8
parent 29 = 20
parent 52 = 30
parent _  = 0

lca     :: Int -> [Int] -> Int
lca w xs | x == y    = x
         | w == 0    = lca y [x, y]
         | y == 0    = lca w [parent x, w]
         | otherwise = lca w [x, parent y]
         where x = head xs
               y = last xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . lca 0 . map read . words) $ lines input
