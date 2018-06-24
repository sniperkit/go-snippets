import System.Environment (getArgs)
import Data.List (intercalate, permutations, sort)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . foldr (\ x xs -> if not (null xs) && head xs == x then xs else x:xs) [] . sort . permutations) $ lines input
