import System.Environment (getArgs)
import Data.Bits (bit)

toTuple       :: [Int] -> (Int, Int)
toTuple [x, y] = (x, y)

haszero    :: Int -> Int -> Bool
haszero 0 0 = True
haszero _ 0 = False
haszero x y | mod y 2 == 1 = haszero x (div y 2)
            | x == 0       = False
            | otherwise    = haszero (pred x) (div y 2)

onezero       :: (Int, Int) -> Int
onezero (x, y) | y < bit x   = 0
               | haszero x y = succ $ onezero (x, pred y)
               | otherwise   = onezero (x, pred y)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . onezero . toTuple . map read . words) $ lines input
