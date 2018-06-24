import System.Environment (getArgs)
import Data.List.Split (splitOn)

fiwr          :: String -> [Int] -> String
fiwr _  []     = ""
fiwr xs (y:ys) = (xs!!(y-1)) : fiwr xs ys

fiwri         :: [String] -> String
fiwri [_]      = ""
fiwri [xs, ys] = fiwr xs (map read (words ys))

nonempty  :: [[a]] -> [[a]]
nonempty s = [x | x <- s, not (null x)]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . nonempty . map (fiwri . splitOn "| ") $ lines input
