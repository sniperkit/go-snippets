import System.Environment (getArgs)
import Data.List (intercalate, nub)
import Data.List.Split (splitOn)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . nub . splitOn ",") $ lines input
