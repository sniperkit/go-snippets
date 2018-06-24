import System.Environment (getArgs)
import Data.List.Split (splitOn)

multiples        :: [Integer] -> Integer
multiples [a, b] | c == 0    = a
                 | otherwise = a - c + b
                 where c = a - floor (fromIntegral a * 2 ** (- logBase 2 (fromIntegral b))) * round (2 ** logBase 2 (fromIntegral b))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . multiples . map read . splitOn ",") $ lines input
