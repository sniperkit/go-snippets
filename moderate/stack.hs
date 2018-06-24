import System.Environment (getArgs)

type Stack = [Int]

push     :: Int -> Stack -> Stack
push x ys = x : ys

pop       :: Stack -> (Int, Stack)
pop (x:xs) = (x, xs)

doStack          :: Stack -> [Int] -> [Int]
doStack []  []    = []
doStack [x] []    = [x]
doStack xs  []    = x : doStack zs []
                  where (x, ys) = pop xs
                        (_, zs) = pop ys
doStack xs (y:ys) = doStack (push y xs) ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map show . doStack [] . map read . words) $ lines input
