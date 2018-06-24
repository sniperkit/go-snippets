import System.Environment (getArgs)

queryCol           :: Int -> [[Int]] -> Int
queryCol _ []       = 0
queryCol x (ys:yss) = ys !! x + queryCol x yss

setCol             :: Int -> Int -> [[Int]] -> [[Int]]
setCol _ _ []       = []
setCol p v (xs:xss) = (take p xs ++ [v] ++ drop (p+1) xs) : setCol p v xss

queryboard           :: [[Int]] -> [String] -> [String]
queryboard _   []     = []
queryboard xss (y:ys) | c == "QueryRow" = show (sum (xss !! p)) : queryboard xss ys
                      | c == "QueryCol" = show (queryCol p xss) : queryboard xss ys
                      | c == "SetRow"   = queryboard (take p xss ++ [replicate 256 v] ++ drop (p+1) xss) ys
                      | c == "SetCol"   = queryboard (setCol p v xss) ys
                      | otherwise       = ["Error"]
                      where cs = words y
                            c  = head cs
                            p  = read (cs !! 1)
                            v  = read (last cs)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . queryboard (replicate 256 (replicate 256 0)) $ lines input
