import System.Environment (getArgs)
import Data.List (isInfixOf, isPrefixOf)
import Data.List.Split (splitOn)

concatstr              :: [(String, Bool)] -> String
concatstr []            = []
concatstr ((xs, _):yss) = xs ++ concatstr yss

prefix          :: String -> String -> String
prefix xs []     = []
prefix xs (y:ys) | isPrefixOf xs (y:ys) = []
                 | otherwise            = y : prefix xs ys

strsu               :: [(String, Bool)] -> String -> String -> [String] -> Int -> String
strsu xs ys zs ass i | i == length xs                       = strsub xs ass
                     | snd (xs!!i) || not (isInfixOf ys fs) = strsu xs ys zs ass (succ i)
                     | fs == ys                             = strsu (take i xs ++ ((zs, True) : drop (i+1) xs)) ys zs ass (succ i)
                     | isPrefixOf ys fs                     = strsu (take i xs ++ ((zs, True) : (drop (length ys) fs, False) : drop (i+1) xs)) ys zs ass (succ i)
                     | otherwise                            = strsu (take i xs ++ ((ps, False) : (zs, True) : (drop (length ps + length ys) fs, False) : drop (i+1) xs)) ys zs ass (i+2)
                     where fs = fst (xs!!i)
                           ps = prefix ys fs

strsub               :: [(String, Bool)] -> [String] -> String
strsub xs []          = concatstr xs
strsub xs (ys:zs:ass) = strsu xs ys zs ass 0

strsubs         :: [String] -> String
strsubs [xs, ys] = strsub [(xs, False)] (splitOn "," ys)

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (strsubs . splitOn ";") $ lines input
