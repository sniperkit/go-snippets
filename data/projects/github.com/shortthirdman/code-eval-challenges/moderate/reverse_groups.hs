import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate)

revgrou     :: [Int] -> Int -> [Int]
revgrou xs y | length xs < y = xs
             | otherwise     = reverse (take y xs) ++ revgrou (drop y xs) y

revgroup         :: [String] -> [Int]
revgroup [xs, ys] = revgrou (map read $ splitOn "," xs) (read ys)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . map show . revgroup . splitOn ";") $ lines input
