import System.Environment (getArgs)

passdown                  :: [Int] -> [Int] -> [Int]
passdown [x] [y, z]        = [y, x + z]
passdown (w:x:xs) (y:z:zs) = y : passdown (x:xs) (max w x + z : zs)

passtriangle     :: [[Int]] -> Int
passtriangle [xs] = maximum xs
passtriangle xss  = passtriangle (passdown xs (head xs + head ys : tail ys) : zss)
                  where xs  = head xss
                        ys  = head (tail xss)
                        zss = tail (tail xss)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    print . passtriangle . map (map read . words) $ lines input
