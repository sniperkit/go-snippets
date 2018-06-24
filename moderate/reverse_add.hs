import System.Environment (getArgs)

pali  :: Integer -> Bool
pali x = x == reve x

reve  :: Integer -> Integer
reve x = read . reverse $ show x

revadd    :: Integer -> Integer -> String
revadd x y | pali y    = show x ++ " " ++ show y
           | x > 100   = "not found"
           | otherwise = revadd (x + 1) (y + reve y)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (revadd 0 . read) $ lines input
