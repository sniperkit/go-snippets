import System.Environment (getArgs)
import Data.List (elemIndices)

replaceTrack       :: Int -> Char -> String -> String
replaceTrack i c xs = take i xs ++ [c] ++ drop (i+1) xs

race               :: Int -> [String] -> [String] -> [String]
race _    xs []     = xs
race (-1) xs (y:ys) = race t [replaceTrack t '|' y] ys
                    where t  = head (elemIndices 'C' y ++ elemIndices '_' y)
race i    xs (y:ys) = xs ++ race t [replaceTrack t sy y] ys
                    where t  = head ([x | x <- elemIndices 'C' y, abs (x - i) < 2] ++
                                     [x | x <- elemIndices '_' y, abs (x - i) < 2])
                          sy | t > i     = '\\'
                             | t < i     = '/'
                             | otherwise = '|'

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . race (-1) [] $ lines input
