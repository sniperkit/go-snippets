import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate, transpose)

slide            :: Int -> [Int] -> [Int]
slide (-1) (y:ys) = slide y ys
slide x []        = [x]
slide x (y:ys)    | y == 0    = slide x ys ++ [0]
                  | x == 0    = slide y ys ++ [0]
                  | x == y    = (x + y) : slide 0 ys
                  | otherwise = x : slide y ys

m2t    :: [[Int]] -> String
m2t xss = intercalate "|" (map (unwords . map show) xss)

tfe :: String -> Int -> [[Int]] -> [[Int]]
tfe "LEFT"  _ xss = map (slide (-1)) xss
tfe "RIGHT" _ xss = map (reverse . slide (-1) . reverse) xss
tfe "UP"    _ xss = transpose $ map (slide (-1)) (transpose xss)
tfe "DOWN"  _ xss = transpose $ map (reverse . slide (-1) . reverse) (transpose xss)

twentyFortyEight             :: [String] -> String
twentyFortyEight [xs, ys, zs] = m2t $ tfe xs (read ys) (map (map read . words) $ splitOn "|" zs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (twentyFortyEight . splitOn "; ") $ lines input
