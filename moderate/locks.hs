import System.Environment (getArgs)

r1 :: [Int]
r1 = [1, 0, 0, 0, 1, 1]
r2 :: [Int]
r2 = [1, 0, 1, 0, 1, 1]

bnot :: Int -> Int
bnot 0 = 1
bnot _ = 0

pmod6  :: Int -> Int
pmod6 x | mod x 6 == 0 = 6
        | otherwise    = mod x 6

locks :: [Int] -> Int
locks [0, _] = 0
locks [x, 0] = x
locks [x, 1] = x-1
locks [x, y] | x > 6 && mod y 2 == 0 = div x 6 * 3 + locks [pmod6 x, y]
             | x > 6                 = div x 6 * 4 + locks [pmod6 x, y]
             | mod y 2 == 0          = sum (init xs) + bnot (last xs)
             | otherwise             = sum (init ys) + bnot (last ys)
             where xs = take x r1
                   ys = take x r2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . locks . map read . words) $ lines input
