import System.Environment (getArgs)
import Data.Char (isAlpha, toLower, toUpper)

rollerCoaster             :: Bool -> String -> String -> String
rollerCoaster _  xs []     = xs
rollerCoaster up xs (y:ys) | isAlpha y && up = rollerCoaster False (xs ++ [toUpper y]) ys
                           | isAlpha y       = rollerCoaster True (xs ++ [toLower y]) ys
                           | otherwise       = rollerCoaster up (xs ++ [y]) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (rollerCoaster True "") $ lines input
