import System.Environment (getArgs)

repetition       :: Int -> Int -> String -> String
repetition x y xs | x == length xs     = show y
                  | xs!!x /= xs!!(x-y) = repetition (succ x) (succ x) xs
                  | otherwise          = repetition (succ x) y xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (repetition 1 1) $ lines input
