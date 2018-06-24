import System.Environment (getArgs)
import Numeric (showIntAtBase)
import Data.Char (intToDigit)

binary  :: Integer -> String
binary i = showIntAtBase 2 intToDigit i ""

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (binary . read) $ [x | x <- lines input, not (null x)]
