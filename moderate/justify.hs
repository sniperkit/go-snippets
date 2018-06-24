import System.Environment (getArgs)

spaces       :: [String] -> Int -> Int -> [String]
spaces xs 0 y = xs
spaces xs x y | x >= y    = spaces (map (++ " ") (init xs) ++ [last xs]) (x-y) y
              | otherwise = map (++ " ") (take x xs) ++ drop x xs

justi :: String -> String
justi xs | length xs == 80 = xs
         | otherwise       = concat (spaces ws l (length ws - 1))
         where ws = words xs
               l  = 80 - sum (map length ws)

maxwords :: String -> Int
maxwords [] = -1
maxwords xs | length xs > 81 = maxwords (take 81 xs)
            | last xs == ' ' = length xs - 1
            | otherwise      = maxwords (init xs)

justif   :: String -> [String]
justif xs | head xs == ' '  = justif (tail xs)
          | length xs <= 80 = [xs]
          | otherwise       = justi (take m xs) : justif (drop m xs)
          where m = maxwords xs

justify   :: String -> String
justify xs | length xs <= 80 = xs
           | otherwise       = init (unlines (justif xs))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map justify $ lines input
