import System.Environment (getArgs)

dsig  :: Int -> Int
dsig 0 = 0
dsig x | mod x 10 > 0 = 2^(3 * mod x 10) + dsig (div x 10)
       | otherwise    = dsig (div x 10)

nextNu    :: Int -> Int -> Int
nextNu x y | x == dsig y = y
           | otherwise   = nextNu x (y + 9)

nextNum  :: Int -> Int
nextNum x = nextNu (dsig x) (x + 9)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . nextNum . read) $ lines input
