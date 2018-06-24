import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.Bits (testBit)

compareBits          :: [Int] -> String
compareBits [i, a, b] | testBit i (a - 1) == testBit i (b - 1) = "true"
                      | otherwise                              = "false"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (compareBits . map read . splitOn ",") $ lines input
