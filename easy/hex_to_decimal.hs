import System.Environment (getArgs)
import Numeric (readHex)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . fst . head . readHex) $ lines input
