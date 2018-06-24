import System.Environment (getArgs)

count :: (Eq a) => a -> [a] -> Int
count x ys = length (filter (== x) ys)

nrc        :: Int -> String -> String -> String
nrc i xs ys | i == length ys = ""
            | elem y xs      = nrc (succ i) xs ys
            | c > 1          = nrc (succ i) (y:xs) ys
            | otherwise      = [y]
            where y = ys!!i
                  c = count y ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (nrc 0 "") $ lines input
