import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate, sort, transpose)

sortmatri    :: [[String]] -> [[Int]]
sortmatri xss = transpose (sort (transpose (map (map read) xss)))

sortmatrix   :: [String] -> String
sortmatrix xs = intercalate " | " (map (unwords . map show) (sortmatri (map words xs)))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (sortmatrix . splitOn " | ") $ lines input
