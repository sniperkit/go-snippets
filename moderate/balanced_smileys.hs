import System.Environment (getArgs)

isBalanced         :: Int -> String -> Bool
isBalanced 0 []     = True
isBalanced _ []     = False
isBalanced (-1) _   = False
isBalanced c (x:xs) | notElem x ":()"            = isBalanced c xs
                    | notElem (last (x:xs)) "()" = isBalanced c (x : init xs)
                    | x == '('                   = isBalanced (succ c) xs
                    | x == ')'                   = isBalanced (pred c) xs
                    | x == ':' && head xs == '(' = isBalanced c (tail xs) || isBalanced (succ c) (tail xs)
                    | x == ':' && head xs == ')' = isBalanced c (tail xs) || isBalanced (pred c) (tail xs)
                    | x == ':'                   = isBalanced c xs
                    | otherwise                  = False

balance   :: String -> String
balance xs | isBalanced 0 xs = "YES"
           | otherwise       = "NO"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map balance $ lines input
