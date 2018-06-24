import System.Environment (getArgs)
import Data.List.Split (splitOn)

complete          :: [Int] -> [String] -> [String]
complete (x:xs) ys | elem z ys = complete xs ys
                   | otherwise = ys ++ [z]
                   where z = show x

recovery                :: [String] -> [[String]] -> [String]
recovery [] [xs, ys]     = recovery (replicate (length xs) []) [xs, complete [1..(length ys + 1)] ys]
recovery zs [[], []]     = zs
recovery zs [x:xs, y:ys] = recovery (take (p - 1) zs ++ [x] ++ drop p zs) [xs, ys]
                             where p = read y

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . recovery [] . map words . splitOn ";") $ lines input
