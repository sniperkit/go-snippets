import System.Environment (getArgs)

pic   :: [Double] -> String
pic xs | x*x + y*y <= r*r = "true"
       | otherwise        = "false"
       where x = head xs - xs!!3
             y = last xs - xs!!1
             r = xs!!2

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (pic . map read . words) $ lines [x | x <- input, elem x "0123456789 -.\n"]
