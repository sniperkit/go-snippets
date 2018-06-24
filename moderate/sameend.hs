import System.Environment (getArgs)
import Data.List (isSuffixOf)
import Data.List.Split (splitOn)

isSuff         :: [String] -> String
isSuff [xs, ys] | isSuffixOf ys xs = "1"
                | otherwise        = "0"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines $ map (isSuff . splitOn ",") [x | x <- lines input, not (null x)]
