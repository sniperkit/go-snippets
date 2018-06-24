import System.Environment (getArgs)
import Data.List.Split (splitOn)

black       :: Int -> [String] -> String
black _ [ys] = ys
black x yss  | n == (-1) = black x (init yss)
             | otherwise = black x (take n yss ++ drop (succ n) yss)
             where n = pred (mod x (length yss))

blackcard         :: [String] -> String
blackcard [xs, ys] = black (read ys) (words xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (blackcard . splitOn " | ") $ lines input
