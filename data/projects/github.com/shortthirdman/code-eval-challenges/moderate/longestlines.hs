import System.Environment (getArgs)
import Data.List (sortBy)

sortLength    :: [a] -> [b] -> Ordering
sortLength a b | length a > length b = LT
               | otherwise           = GT

merge     :: Int -> [String] -> [String]
merge i xs | length xs <= i = sortBy sortLength xs
           | otherwise      = take i (sortBy sortLength xs)

longest            :: Int -> [String] -> [String] -> [String]
longest _ xs []     = xs
longest 0 _  (y:ys) = longest (read y) [] ys
longest i xs (y:ys) = longest i (merge i (y : xs)) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . longest 0 [] $ lines input
