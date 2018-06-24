import System.Environment (getArgs)
import Text.Printf (printf)

angles  :: Double -> String
angles s = printf "%d.%02d'%02d\"" dg mn sc
         where dg = floor s :: Int
               mn = floor ((s - fromIntegral dg) * 60) :: Int
               sc = floor (((s - fromIntegral dg) * 60 - fromIntegral mn) * 60) :: Int

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (angles . (read :: String -> Double)) $ lines input
