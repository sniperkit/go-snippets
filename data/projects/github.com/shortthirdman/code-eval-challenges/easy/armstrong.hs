import System.Environment (getArgs)
import Data.Char (digitToInt)

nthp    :: Int -> Int -> Int
nthp e i = i ^ e

sumofnthp  :: String -> String
sumofnthp i = show . sum . map (nthp e) $ map digitToInt i
            where e = floor (logBase 10 (read i::Double)) + 1

armstrong  :: String -> String
armstrong i | sumofnthp i == i = "True"
            | otherwise        = "False"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map armstrong $ lines input
