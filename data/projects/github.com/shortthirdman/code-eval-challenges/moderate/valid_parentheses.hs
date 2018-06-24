import System.Environment (getArgs)

parens      :: String -> String -> String
parens [] [] = "True"
parens _ []  = "False"
parens xs ys | y == ')' && not (null xs) && x == '(' = parens (tail xs) (tail ys)
             | y == ']' && not (null xs) && x == '[' = parens (tail xs) (tail ys)
             | y == '}' && not (null xs) && x == '{' = parens (tail xs) (tail ys)
             | elem y ")]}"                          = "False"
             | otherwise                             = parens (y : xs) (tail ys)
             where x = head xs
                   y = head ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (parens "") $ lines input
