import System.Environment (getArgs)
import Data.List (elemIndices)

readMore  :: String -> String
readMore s | length s < 56 = s
           | otherwise     = take (last e) s ++ "... <Read More>"
           where e = 40 : [x | x <- elemIndices ' ' s, x < 40]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map readMore $ lines input
