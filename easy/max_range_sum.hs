import System.Environment (getArgs)
import Data.List.Split (splitOn)

maxran                  :: Int -> Int -> [Int] -> [Int] -> Int
maxran x _ _      []     = x
maxran x c (y:ys) (z:zs) = maxran (max x d) d (ys ++ [z]) zs
                         where d = c - y + z

maxrange         :: [String] -> Int
maxrange [ns, xs] = maxran (max 0 z) z zs (drop n ys)
                  where n  = read ns
                        ys = map read (words xs)
                        zs = take n ys
                        z  = sum zs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . maxrange . splitOn ";") $ lines input
