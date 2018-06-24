import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (sort)

absurd       :: [String] -> String
absurd [_]    = ""
absurd (x:xs) | x == head xs = x
              | otherwise    = absurd xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (absurd . sort . splitOn "," . last . splitOn ";") $ lines input
