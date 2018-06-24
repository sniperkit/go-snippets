import System.Environment (getArgs)

uniq     :: (Eq a) => [a] -> a -> Bool
uniq xs y = length (filter (== y) xs) == 1

cycledet          :: String -> [String] -> [String]
cycledet _  []     = []
cycledet [] (y:ys) = y : cycledet y ys
cycledet x  (y:ys) | x == y    = []
                   | otherwise = y : cycledet x ys

cycledetection   :: [String] -> [String]
cycledetection [] = []
cycledetection xs = cycledet [] (dropWhile (uniq xs) xs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . cycledetection . words) $ lines input
