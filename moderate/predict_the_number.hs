import System.Environment (getArgs)
import Data.Bits (popCount)

count  :: Int -> String
count i = show $ mod (popCount i) 3

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (count . (read :: String -> Int)) $ lines input
