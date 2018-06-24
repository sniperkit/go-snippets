import System.Environment (getArgs)

clean  :: Char -> Char
clean c | elem c ['a'..'j'] = toEnum (fromEnum c - 49)
        | otherwise         = c

hidden  :: String -> String
hidden s | cs == ""  = "NONE"
         | otherwise = cs
         where cs = [clean x | x <- s, elem x (['0'..'9'] ++ ['a'..'j'])]

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map hidden $ lines input
