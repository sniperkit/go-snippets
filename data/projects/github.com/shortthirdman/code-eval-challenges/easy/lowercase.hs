import System.Environment (getArgs)
import Data.Char (toLower)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (map toLower) $ lines input
