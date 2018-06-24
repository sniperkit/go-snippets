import System.Environment (getArgs)
import Data.List.Split (splitOn)

subs         :: [String] -> Int
subs [_,  []] = 1
subs [[], _ ] = 0
subs [xs, ys] | head xs == head ys = subs [tail xs, tail ys] + subs [tail xs, ys]
              | otherwise          = subs [tail xs, ys]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . subs . splitOn ",") $ lines input
