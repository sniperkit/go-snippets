import System.Environment (getArgs)
import Data.List (intercalate)

factors  :: Int -> [Int]
factors n = [x | x <- [1..n], mod n x == 0]

prime  :: Int -> Bool
prime x = factors x == [1, x]

m  :: Int -> Int
m x = 2^x - 1

mersenne  :: Int -> [Int]
mersenne n = map m $ takeWhile (\x -> m x < n) [x | x <- 2 : [3, 5..], prime x]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate ", " . map show . mersenne . read) $ lines input
