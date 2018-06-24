import Text.Printf (printf)

table      :: [[String]] -> [String] -> [[String]]
table xs [] = xs
table xs ys = table (xs ++ [take 12 ys]) (drop 12 ys)

main :: IO ()
main = putStr . unlines . map concat . table [] $ map (printf "%4d" :: Int -> String) [x * y | x <- [1 .. 12], y <- [1 .. 12]]
