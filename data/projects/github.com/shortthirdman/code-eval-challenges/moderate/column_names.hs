import System.Environment (getArgs)
import Data.Char (chr, ord)

columnnames  :: Int -> String
columnnames 0 = ""
columnnames x = columnnames (div y 26) ++ [chr (ord 'A' + mod y 26)]
              where y = pred x

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (columnnames . read) $ lines input
