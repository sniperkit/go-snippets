import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (elemIndices)

lastIndex  :: [String] -> String
lastIndex s | null e    = "-1"
            | otherwise = show . last $ e
            where e = elemIndices (head (last s)) (head s)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (lastIndex . splitOn ",") $ [x | x <- lines input, not (null x)]
