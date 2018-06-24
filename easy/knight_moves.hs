import System.Environment (getArgs)

moves :: [(Char -> Char, Char -> Char)]
moves  = [(pred . pred, pred), (pred . pred, succ), (pred, pred . pred), (pred, succ . succ),
          (succ, pred . pred), (succ, succ . succ), (succ . succ, pred), (succ . succ, succ)]

knigh           :: Char -> Char -> [(Char -> Char, Char -> Char)] -> [String]
knigh _ _ []     = []
knigh x y (z:zs) | elem (fst z x) ['a'..'h'] && elem (snd z y) ['1'..'8'] = (fst z x : [snd z y]) : knigh x y zs
                 | otherwise = knigh x y zs

knight   :: String -> [String]
knight [x, y] = knigh x y moves

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . knight) $ lines input
