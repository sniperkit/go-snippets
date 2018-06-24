import System.Environment (getArgs)
import Data.List (intercalate)

tele          :: String -> String -> [String]
tele xs []     = [xs]
tele xs (y:ys) | y == '2'  = tele (xs ++ "a") ys ++ tele (xs ++ "b") ys ++ tele (xs ++ "c") ys
               | y == '3'  = tele (xs ++ "d") ys ++ tele (xs ++ "e") ys ++ tele (xs ++ "f") ys
               | y == '4'  = tele (xs ++ "g") ys ++ tele (xs ++ "h") ys ++ tele (xs ++ "i") ys
               | y == '5'  = tele (xs ++ "j") ys ++ tele (xs ++ "k") ys ++ tele (xs ++ "l") ys
               | y == '6'  = tele (xs ++ "m") ys ++ tele (xs ++ "n") ys ++ tele (xs ++ "o") ys
               | y == '7'  = tele (xs ++ "p") ys ++ tele (xs ++ "q") ys ++ tele (xs ++ "r") ys ++ tele (xs ++ "s") ys
               | y == '8'  = tele (xs ++ "t") ys ++ tele (xs ++ "u") ys ++ tele (xs ++ "v") ys
               | y == '9'  = tele (xs ++ "w") ys ++ tele (xs ++ "x") ys ++ tele (xs ++ "y") ys ++ tele (xs ++ "z") ys
               | otherwise = tele (xs ++ [y]) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . tele "") $ lines input
