import System.Environment (getArgs)
import Data.List (sortBy)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . sortBy (flip compare) . words) $ lines input
