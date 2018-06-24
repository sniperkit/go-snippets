import System.Environment (getArgs)
import Data.List.Split (splitOn)
import Data.List (intercalate, nub, sort)
import Control.Monad (replicateM)

stringlist         :: [String] -> [String]
stringlist [ns, xs] = replicateM (read ns) (nub (sort xs))

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (intercalate "," . stringlist . splitOn ",") $ lines input
