import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate)

fmt0  :: Integer -> String
fmt0 i | i < 10    = '0' : show i
       | otherwise = show i

nice  :: Integer -> String
nice i = intercalate ":" $ map fmt0 t
       where t = div i 3600 : mod (div i 60) 60 : [mod i 60]

delta  :: [[Integer]] -> Integer
delta [[h1, m1, s1], [h2, m2, s2]] = abs (h1*3600 + m1*60 + s1 - h2*3600 - m2*60 - s2)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (nice . delta . map (map read . splitOn ":") . words) $ lines input
