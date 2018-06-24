import System.Environment (getArgs)
import Text.Regex.Posix ((=~))

regex :: String
regex = "^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@[0-9A-Za-z_]*([0-9A-Za-z_]+[.])+[0-9A-Za-z_]{2,4}$"

email   :: String -> String
email xs | xs =~ regex :: Bool = "true"
         | otherwise           = "false"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map email $ lines input
