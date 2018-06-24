import System.Environment (getArgs)

isEven  :: Integer -> Integer
isEven i = mod (succ i) 2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . isEven . read) $ lines input
