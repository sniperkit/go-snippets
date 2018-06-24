import System.Environment (getArgs)

b :: [Integer]
b = [0, 1, 2, 1, 2]

coin   :: Integer -> String
coin i = show (div i 5 + b !! fromInteger (rem i 5))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (coin . (read :: String -> Integer)) $ lines input
