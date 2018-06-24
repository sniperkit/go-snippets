import System.Environment (getArgs)
import Data.List.Split (splitOn)

rotate           :: Int -> [String] -> String
rotate x [xs, ys] | x == length xs = "False"
                  | xs == ys       = "True"
                  | otherwise      = rotate (succ x) [last xs : init xs, ys]
rotate _ _        = "False"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (rotate 0 . splitOn ",") $ lines input
