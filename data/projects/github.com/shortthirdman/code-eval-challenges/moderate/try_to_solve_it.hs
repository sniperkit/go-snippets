import System.Environment (getArgs)
import Data.Char (chr, ord)

trytosolve  :: Char -> Char
trytosolve x | x >= 'a' && x <= 'f' = chr (ord x + 20)
             | x >= 'g' && x <= 'm' = chr (ord x + 7)
             | x >= 'n' && x <= 't' = chr (ord x - 7)
             | x >= 'u' && x <= 'z' = chr (ord x - 20)
             | otherwise            = x

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (map trytosolve) $ lines input
