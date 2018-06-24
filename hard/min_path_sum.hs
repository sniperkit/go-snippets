import System.Environment (getArgs)
import Data.List.Split (splitOn)

infinity :: Int
infinity  = maxBound :: Int

minpa                    :: [Int] -> [Int] -> [Int]
minpa _         [y]       = [y]
minpa (_:x':xs) (y:y':ys) = y : minpa (x':xs) ((min y x' + y') : ys)

minpat            :: [Int] -> [[Int]] -> String
minpat xs []       = show (last xs)
minpat xs (ys:yss) = minpat (minpa xs ((head xs + head ys) : tail ys)) yss

minpath         :: [String] -> [String]
minpath []       = []
minpath (xs:xss) = minpat (0 : replicate (n-1) infinity) (map (map read . splitOn ",") (take n xss)) : minpath (drop n xss)
                 where n = read xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . minpath $ lines input
