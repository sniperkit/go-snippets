import System.Environment (getArgs)
import Data.Char (isAlpha, toLower)

cleanWord          :: String -> String -> [String]
cleanWord xs []     = [xs]
cleanWord xs (y:ys) | isAlpha y = cleanWord (xs ++ [toLower y]) ys
                    | otherwise = xs : cleanUp ys

cleanUp       :: String -> [String]
cleanUp []     = []
cleanUp (x:xs) | not (isAlpha x) = cleanUp xs
               | otherwise       = cleanWord [toLower x] xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . cleanUp) $ lines input
