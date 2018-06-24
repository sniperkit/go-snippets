import System.Environment (getArgs)
import Data.List (isInfixOf)

findrep     :: Int -> String -> String
findrep x ys | 2 * x > length ys                 = ""
             | take x ys == replicate x ' '      = findrep x (tail ys)
             | isInfixOf (take x ys) (drop x ys) = take x ys
             | otherwise                         = findrep x (tail ys)

repeated      :: String -> String -> String
repeated xs ys | null (findrep (length xs + 1) ys) = xs
               | otherwise                         = repeated (findrep (length xs + 1) ys) ys

repeatedsub   :: String -> String
repeatedsub xs | null (repeated "" xs) = "NONE"
               | otherwise             = repeated "" xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map repeatedsub $ lines input
