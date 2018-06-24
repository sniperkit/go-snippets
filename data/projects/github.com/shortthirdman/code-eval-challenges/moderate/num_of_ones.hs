import System.Environment (getArgs)
import Data.Bits (popCount)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . popCount . (read :: String -> Int)) $ lines input
