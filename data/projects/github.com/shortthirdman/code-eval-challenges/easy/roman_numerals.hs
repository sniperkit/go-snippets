import System.Environment (getArgs)

ronum :: [Int]
ronum = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]
rostr :: [String]
rostr = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"]

roman      :: String -> Int -> Int -> String
roman s _ 0 = s
roman s j i | ronum !! j > i = roman s (succ j) i
            | otherwise      = roman (s ++ rostr !! j) j (i - ronum !! j)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (roman "" 0 . read) $ lines input
