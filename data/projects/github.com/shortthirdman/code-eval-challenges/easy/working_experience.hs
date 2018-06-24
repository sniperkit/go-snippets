import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (nub)

moy  :: String -> Int
moy s | s == "Jan" = 0
      | s == "Feb" = 1
      | s == "Mar" = 2
      | s == "Apr" = 3
      | s == "May" = 4
      | s == "Jun" = 5
      | s == "Jul" = 6
      | s == "Aug" = 7
      | s == "Sep" = 8
      | s == "Oct" = 9
      | s == "Nov" = 10
      | otherwise  = 11

month  :: String -> Int
month s = 12 * (k - 1990) + moy (take 3 s)
        where k = read (drop 4 s)

workm       :: [String] -> [Int]
workm []     = []
workm (x:xs) = [(month y)..(month z)] ++ workm xs
             where [y, z] = splitOn "-" x

work  :: [String] -> String
work s = show $ div (length . nub $ workm s) 12

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (work . splitOn "; ") $ lines input
