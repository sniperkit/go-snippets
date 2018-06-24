import System.Environment (getArgs)
import Data.List (sort)

bubbl       :: Int -> Int -> [Int] -> [Int]
bubbl x y zs | x >= length zs - y = zs
             | otherwise          = bubbl (x+1) y (take x zs ++ sort (take (y+1) (drop x zs)) ++ drop (x+y+1) zs)

bubble :: [String] -> [Int]
bubble xs = bubbl 0 (min (read (last xs)) (length xs - 3)) (map read (init (init xs)))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . bubble . words) $ lines input
