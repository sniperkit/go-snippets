import System.Environment (getArgs)

age  :: Int -> String
age i | i < 0 || i > 100 = "This program is for humans"
      | i < 3            = "Still in Mama's arms"
      | i < 5            = "Preschool Maniac"
      | i < 12           = "Elementary school"
      | i < 15           = "Middle school"
      | i < 19           = "High school"
      | i < 23           = "College"
      | i < 66           = "Working for the man"
      | otherwise        = "The Golden Years"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (age . read) $ lines input
