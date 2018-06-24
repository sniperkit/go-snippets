import System.Environment (getArgs)
import Data.List (sort)
import Numeric (showFFloat)

formatFloat3  :: Double -> String
formatFloat3 f = showFFloat (Just 3) f ""

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (unwords . map formatFloat3 . sort . map (read :: String -> Double) . words) $ lines input
