import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.Ratio (Ratio, (%))

clean   :: String -> String
clean xs = [x | x <- xs, notElem x "()"]

cityBlock                       :: Int -> [Ratio Int] -> [Ratio Int] -> Int
cityBlock intersec []     _      = intersec + 1
cityBlock intersec _      []     = intersec + 1
cityBlock intersec (x:xs) (y:ys) | x == y    = cityBlock (succ intersec) xs ys
                                 | x < y     = cityBlock intersec xs (y:ys)
                                 | otherwise = cityBlock intersec (x:xs) ys

cityBlocks         :: [String] -> Int
cityBlocks [xs, ys] | length ps < 2 || length qs < 2 = 0
                    | otherwise                      = length ps + length qs - cityBlock 0 ps'' qs''
                    where ps   = map (read :: String -> Int) (splitOn "," xs)
                          qs   = map (read :: String -> Int) (splitOn "," ys)
                          ps'  = map (+ (-(head ps))) ps
                          qs'  = map (+ (-(head qs))) qs
                          stl  = last ps'
                          avl  = last qs'
                          ps'' = map fromIntegral ps'
                          qs'' = map (\x -> (x * stl) % avl) qs'

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (show . cityBlocks . words) . lines $ clean input
