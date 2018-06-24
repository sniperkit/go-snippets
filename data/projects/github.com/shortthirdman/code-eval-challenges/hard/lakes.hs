import System.Environment (getArgs)
import Data.Bits ((.|.))
import Data.List.Split (splitOn)

lak           :: Int -> Int -> [Int] -> [Int]
lak _ _ []     = []
lak x y (z:zs) | z == x    = y : lak x y zs
               | otherwise = z : lak x y zs

lake                         :: Int -> [Int] -> [String] -> [[String]] -> Int
lake i xs     []     []       = 0
lake i xs     []     (zs:zss) = lake i (tail xs ++ [0]) zs zss
lake i (x:xs) (y:ys) zss      | y == "#"                     = lake i (xs ++ [0]) ys zss
                              | ll > 0 && rl > 0 && ll /= rl = (-1) + lake i (lak rl ll xs ++ [ll]) ys zss
                              | ll > 0                       = lake i (xs ++ [ll]) ys zss
                              | rl > 0                       = lake i (xs ++ [rl]) ys zss
                              | otherwise                    = 1 + lake (succ i) (xs ++ [i]) ys zss
                              where ll = (.|.) x ((.|.) (head xs) (last xs))
                                    rl = xs!!1

lakes         :: [[String]] -> Int
lakes (xs:xss) = lake 1 (replicate (length xs + 2) 0) xs xss

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . lakes . map words . splitOn "|") $ lines input
