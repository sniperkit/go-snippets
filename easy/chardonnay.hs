import System.Environment (getArgs)
import Data.List.Split (splitOn)

findWine          :: String -> String -> Bool
findWine _  ""     = True
findWine xs (y:ys) | elem y xs = findWine (l ++ r) ys
                   | otherwise = False
                   where l = takeWhile (/= y) xs
                         r = tail (dropWhile (/= y) xs)

findWines            :: [String] -> String -> [String]
findWines []       _  = []
findWines (xs:xss) ys | findWine xs ys = xs : findWines xss ys
                      | otherwise      = findWines xss ys

findAllWines         :: [String] -> String
findAllWines [xs, ys] | null w    = "False"
                      | otherwise = unwords w
                      where w = findWines (words xs) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (findAllWines . splitOn " | ") $ lines input
