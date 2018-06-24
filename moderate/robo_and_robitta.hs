import System.Environment (getArgs)
import Data.List.Split (splitOn)

robi        :: Int -> Int -> Int -> Int -> Int
robi a b x y | b == y    = x
             | otherwise = a + robi (pred b) a (b - y) x

robo         :: [String] -> String
robo [xs, ys] = show $ robi (head as) (last as) (head bs) (last bs)
              where as = map read $ splitOn "x" xs
                    bs = map read $ words ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (robo . splitOn "|") $ lines input
