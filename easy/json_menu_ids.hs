import System.Environment (getArgs)
import Data.Char (isDigit)
import Data.List (isPrefixOf)

jsonMenu     :: Int -> String -> Int
jsonMenu a [] = a
jsonMenu a xs | not (isPrefixOf "\"id\": " xs) = jsonMenu a (tail xs)
              | head xs'' == ','               = jsonMenu (a+y) xs''
              | otherwise                      = jsonMenu a xs''
              where xs'  = drop 6 xs
                    y    = read (takeWhile isDigit xs')
                    xs'' = dropWhile isDigit xs'

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines $ map (show . jsonMenu 0) [x | x <- lines input, not (null x)]
