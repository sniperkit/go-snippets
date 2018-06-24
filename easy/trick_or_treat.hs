import System.Environment (getArgs)
import Data.List.Split (splitOn)

trick             :: [Int] -> Int
trick [v, z, w, h] | v + z + w == 0 = 0
                   | otherwise      = div ((v * 3 + z * 4 + w * 5) * h) (v + z + w)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . trick . map read . splitOn ",") $ lines [x | x <- input, elem x "0123456789,\n"]
