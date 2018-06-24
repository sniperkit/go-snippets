import System.Environment (getArgs)

m :: String
m = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"

morse            :: Int -> String -> String -> String
morse i xs []     = xs ++ [m !! (i - 2)]
morse i xs (y:ys) | y == '.'  = morse (i * 2) xs ys
                  | y == '-'  = morse (i * 2 + 1) xs ys
                  | i == 1    = morse i (xs ++ " ") ys
                  | otherwise = morse 1 (xs ++ [m !! (i - 2)]) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (morse 1 "") $ lines input
