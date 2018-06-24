import System.Environment (getArgs)

phrases :: [String]
phrases = [", yeah!", ", this is crazy, I tell ya.",
    ", can U believe this?", ", eh?", ", aw yea.",
    ", yo.", "? No way!", ". Awesome!"]

slang           :: Int -> Bool -> String -> String
slang _ _ []     = ""
slang i l (x:xs) | elem x ".!?" && l = phrases!!i ++ slang (mod (i+1) 8) False xs
                 | elem x ".!?"      = x : slang i True xs
                 | otherwise         = x : slang i l xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr $ slang 0 False (if last input == '\n' then input else input ++ "\n")
