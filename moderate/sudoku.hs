import System.Environment (getArgs)
import Data.List (sort, transpose)
import Data.List.Split (splitOn)

unflatten     :: Int -> [Int] -> [[Int]]
unflatten _ [] = []
unflatten x ys = take x ys : unflatten x (drop x ys)

sudosqus     :: Int -> [Int] -> Bool
sudosqus 4 ys = sudorows 4 (take 2 ys ++ take 2 (drop 4 ys))
sudosqus 9 ys = sudorows 9 (take 3 ys ++ take 3 (drop 9 ys) ++ take 3 (drop 18 ys) ++
                            take 3 (drop 3 ys) ++ take 3 (drop 12 ys) ++ take 3 (drop 21 ys) ++
                            take 3 (drop 27 ys) ++ take 3 (drop 36 ys) ++ take 3 (drop 45 ys) ++
                            take 3 (drop 30 ys) ++ take 3 (drop 39 ys) ++ take 3 (drop 48 ys))
sudosqus _ _  = False

sudocols     :: Int -> [Int] -> Bool
sudocols x ys = sudorows x (concat (transpose (unflatten x ys)))

sudorows     :: Int -> [Int] -> Bool
sudorows _ [] = True
sudorows x ys | sort (take x ys) == [1..x] = sudorows x (drop x ys)
              | otherwise                  = False

sudok     :: Int -> [Int] -> Bool
sudok x ys = sudorows x ys && sudocols x ys && sudosqus x ys

sudoku         :: [String] -> String
sudoku [xs, ys] | sudok (read xs) (map read $ splitOn "," ys) = "True"
                | otherwise                                   = "False"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (sudoku . splitOn ";") $ lines input
