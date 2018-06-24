import System.Environment (getArgs)
import Data.List (elem)
import Data.Map (Map)
import qualified Data.Map as Map

codel :: String
codel  = "nug\
    \ rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp\
    \ tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi\
    \ de kr kd eoya kw aej icfkici re zjkr"
decol :: String
decol  = "bjv\
    \ the public is amazed by the quickness of the juggler\
    \ we think that our language is impossible to understand\
    \ so it is okay if you decided to quit"

codeMap                    :: String -> String -> String -> String -> Map Char Char
codeMap (x:xs) (y:ys) cs ds | elem x cs = codeMap xs (y:ys) cs ds
                            | elem y ds = codeMap (x:xs) ys cs ds
                            | otherwise = Map.insert x y (Map.fromList (zip cs ds))

lostInTranslation       :: String -> String
lostInTranslation []     = []
lostInTranslation (x:xs) = Map.findWithDefault ' ' x cMap : lostInTranslation xs
                         where cMap = codeMap ['a'..'z'] ['a'..'z'] codel decol

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map lostInTranslation $ lines input
