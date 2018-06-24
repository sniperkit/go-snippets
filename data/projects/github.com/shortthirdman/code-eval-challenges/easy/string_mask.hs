import System.Environment (getArgs)
import Data.Char (toUpper)

stringmask             :: [String] -> String
stringmask [[], []]     = []
stringmask [x:xs, y:ys] | y == '0'  = x : stringmask [xs, ys]
                        | otherwise = toUpper x : stringmask [xs, ys]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (stringmask . words) $ lines input
