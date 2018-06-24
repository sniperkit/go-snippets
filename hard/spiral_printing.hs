import System.Environment (getArgs)
import Data.List.Split (splitOn)

unflatten     :: Int -> [a] -> [[a]]
unflatten _ [] = []
unflatten x ys = take x ys : unflatten x (drop x ys)

spira         :: [[String]] -> [String]
spira []       = []
spira [xs]     = xs
spira [xs, ys] = xs ++ reverse ys
spira xss      | length xs == 1 = concat xss
               | length xs == 2 = xs ++ map last yss ++ reverse (map head yss)
               | otherwise      = xs ++ map last yss
                                     ++ reverse (last zss)
                                     ++ reverse (map head ass)
                                     ++ spira bss
               where xs  = head xss
                     yss = tail xss
                     zss = map init yss
                     ass = init zss
                     bss = map tail ass

spiral             :: [String] -> [String]
spiral [_, ys, zs] = spira (unflatten (read ys) (words zs))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . spiral . splitOn ";") $ lines input
