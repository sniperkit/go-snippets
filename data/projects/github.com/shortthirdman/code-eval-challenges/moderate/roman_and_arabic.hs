import System.Environment (getArgs)
import Data.Char (digitToInt)

roma :: Char -> Int
roma 'M' = 1000
roma 'D' =  500
roma 'C' =  100
roma 'L' =   50
roma 'X' =   10
roma 'V' =    5
roma _   =    1

romar             :: Int -> Int -> String -> Int
romar x y []       = x * y
romar x y (z:a:zs) | roma a > y = rest - x * y
                   | otherwise  = rest + x * y
                   where rest = romar (digitToInt z) (roma a) zs

romara         :: String -> Int
romara []       = 0
romara (x:y:xs) = romar (digitToInt x) (roma y) xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . romara) $ lines input
