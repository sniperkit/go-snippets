import System.Environment (getArgs)
import Data.Bits (shiftL, shiftR)
import Data.Char (chr, ord)
import Data.List (elemIndex, sort)
import Data.Maybe (fromJust)

bw       :: [Int] -> Int -> Int -> String
bw _  _ 0 = []
bw xs y z = chr (shiftR (xs!!y) 16) : bw xs (mod (xs!!y) (shiftL 1 16)) (z - 1)

mapIndex         :: Int -> [Int] -> [Int]
mapIndex _ []     = []
mapIndex x (y:ys) = (shiftL y 16 + x) : mapIndex (succ x) ys

burrowsWheeler   :: String -> String
burrowsWheeler xs = bw ys k (length xs)
                  where k  = fromJust $ elemIndex '$' xs
                        ys = sort (mapIndex 0 (map ord xs))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (burrowsWheeler . init) $ lines input
