import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate, sort)

clean   :: String -> String
clean xs = [x | x <- xs, notElem x "()[]"]

pileofbric            :: [Int] -> [[String]] -> [String]
pileofbric _  []       = []
pileofbric hs (xs:xss) | head b <= head hs && b!!1 <= hs!!1 = head xs : pileofbric hs xss -- || diagonal b hs
                       | otherwise                          = pileofbric hs xss
                       where t = map read (splitOn "," (xs!!1))
                             u = map read (splitOn "," (xs!!2))
                             b = sort [abs (head t - head u), abs (t!!1 - u!!1), abs (t!!2 - u!!2)]

pileofbrick             :: [String] -> [String] -> [String]
pileofbrick [xs, ys] zss = pileofbric h (map words zss)
                         where s = map read (splitOn "," xs)
                               t = map read (splitOn "," ys)
                               h = sort [abs (head s - head t), abs (s!!1 - t!!1)]

pileofbricks         :: [String] -> String
pileofbricks [xs, ys] | null p    = "-"
                      | otherwise = intercalate "," r
                      where p = pileofbrick (words xs) (splitOn ";" ys)
                            r = map show (sort (map (read :: String -> Int) p))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (pileofbricks . splitOn "|") . lines $ clean input
