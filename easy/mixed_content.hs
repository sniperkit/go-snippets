import System.Environment (getArgs)
import Data.List (intercalate)
import Data.List.Split (splitOn)
import Data.Char (isDigit)

areDigits  :: String -> Bool
areDigits s | length s == 1 && isDigit (head s)      = True
            | length s < 2 || not (isDigit (head s)) = False
            | otherwise                              = areDigits (tail s)

mixed                :: [[String]] -> [String] -> String
mixed [[], ys] []     = intercalate "," ys
mixed [xs, []] []     = intercalate "," xs
mixed [xs, ys] []     = intercalate "," ys ++ "|" ++ intercalate "," xs
mixed [xs, ys] (z:zs) | areDigits z = mixed [xs ++ [z], ys] zs
                      | otherwise   = mixed [xs, ys ++ [z]] zs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (mixed [[], []] . splitOn ",") $ lines input
