import System.Environment (getArgs)
import Data.List.Split (splitOn)

countBugs             :: [String] -> Int
countBugs [[], []]     = 0
countBugs [x:xs, y:ys] | x /= y    = succ (countBugs [xs, ys])
                       | otherwise = countBugs [xs, ys]

prioFix  :: Int -> String
prioFix 0 = "Done"
prioFix x | x <= 2    = "Low"
          | x <= 4    = "Medium"
          | x <= 6    = "High"
          | otherwise = "Critical"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (prioFix . countBugs . splitOn " | ") $ lines input
