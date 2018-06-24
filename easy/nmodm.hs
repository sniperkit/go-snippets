import System.Environment (getArgs)
import Data.List.Split (splitOn)

modulo  :: [Integer] -> Integer
modulo [n, m] = n - div n m * m

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . modulo . map read . splitOn ",") $ lines input
