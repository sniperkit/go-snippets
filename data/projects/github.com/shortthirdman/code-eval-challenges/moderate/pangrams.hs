import System.Environment (getArgs)
import Data.Char (toLower)
import Control.Applicative
import qualified Data.ByteString.Char8 as B

readFileAscii     :: String -> IO String
readFileAscii path = B.unpack <$> B.map (clearChar '-') <$> B.readFile path
    where clearChar    :: Char -> Char -> Char
          clearChar d c | c < '\128' = c
                        | otherwise  = d

pang          :: String -> String -> String
pang []     _  = ""
pang (x:xs) ys | elem x ys = pang xs ys
               | otherwise = x : pang xs ys

pangram  :: String -> String
pangram s | null ls   = "NULL"
          | otherwise = ls
          where ls = pang ['a'..'z'] s

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFileAscii inpFile
    putStr . unlines . map pangram $ lines [toLower x | x <- input, x == '\n' || elem x ['a'..'z'] || elem x ['A'..'Z']]
