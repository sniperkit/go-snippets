import System.Environment (getArgs)

cardval :: Char -> Int
cardval '2' = 2
cardval '3' = 3
cardval '4' = 4
cardval '5' = 5
cardval '6' = 6
cardval '7' = 7
cardval '8' = 8
cardval '9' = 9
cardval '1' = 10
cardval 'J' = 11
cardval 'Q' = 12
cardval 'K' = 13
cardval 'A' = 14
cardval _   = -1

trump                 :: [String] -> String
trump [xs, ys, _, [z]] | xc == z && yc /= z    = xs
                       | xc /= z && yc == z    = ys
                       | cardval x > cardval y = xs
                       | cardval x < cardval y = ys
                       | otherwise             = xs ++ (' ' : ys)
                       where xc = last xs
                             yc = last ys
                             x  = head xs
                             y  = head ys

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (trump . words) $ lines input
