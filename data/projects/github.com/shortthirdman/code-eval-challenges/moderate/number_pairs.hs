import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate)

numpairs :: [Int] -> Int -> [String]
numpairs []     _ = []
numpairs [_]    _ = []
numpairs (x:xs) y | x + z < y = numpairs xs y
                  | x + z > y = numpairs (x : init xs) y
                  | otherwise = (show x ++ "," ++ show z) : numpairs (init xs) y
                  where z = last xs

numberpairs         :: [String] -> String
numberpairs [xs, ys] | null zs   = "NULL"
                     | otherwise = intercalate ";" zs
                     where zs = numpairs (map read $ splitOn "," xs) (read ys)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (numberpairs . splitOn ";") $ lines input
