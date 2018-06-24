import System.Environment (getArgs)
import Data.Char (toUpper)

capitalize         :: String -> String
capitalize []       = []
capitalize (x : xs) = toUpper x : xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map capitalize . words) $ lines input
