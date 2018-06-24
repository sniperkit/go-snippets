import System.Environment (getArgs)
import Data.Char (digitToInt)

digits :: [String]
digits  = [
    "-**----*--***--***---*---****--**--****--**---**--",
    "*--*--**-----*----*-*--*-*----*-------*-*--*-*--*-",
    "*--*---*---**---**--****-***--***----*---**---***-",
    "*--*---*--*-------*----*----*-*--*--*---*--*----*-",
    "-**---***-****-***-----*-***---**---*----**---**--",
    "--------------------------------------------------"]

bigdigit     :: [String] -> Int -> [String]
bigdigit [x0s, x1s, x2s, x3s, x4s, x5s] y = [x0s ++ take 5 (drop (5 * y) (head digits)), x1s ++ take 5 (drop (5 * y) (digits!!1)), x2s ++ take 5 (drop (5 * y) (digits!!2)), x3s ++ take 5 (drop (5 * y) (digits!!3)), x4s ++ take 5 (drop (5 * y) (digits!!4)), x5s ++ take 5 (drop (5 * y) (digits!!5))]

bigdigits :: [String] -> String -> [String]
bigdigits xs []     = xs
bigdigits xs (y:ys) | elem y "0123456789" = bigdigits zs ys
                    | otherwise = bigdigits xs ys
                    where zs = bigdigit xs (digitToInt y)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . concatMap (unlines . bigdigits ["", "", "", "", "", ""]) $ lines input
