import System.Environment (getArgs)

nstr :: [String]
nstr = ["Zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"]

wrd :: [Int] -> String
wrd [0, x, y] = wrd [x, y]
wrd [0, x]    = wrd [x]
wrd [0]       = ""
wrd [x, y, z] = nstr!!x ++ "Hundred" ++ wrd [y, z]
wrd [x, y]    | x > 1     = nstr!!(x + 18) ++ wrd [y]
              | otherwise = nstr!!(10 + y)
wrd [x]       = nstr!!x

dolla                            :: [Int] -> String
dolla [0, 0, 0, x, y, z, a, b, c] = dolla [x, y, z, a, b, c]
dolla [0, 0, 0, x, y, z]          = dolla [x, y, z]
dolla [0, 0, 0]                   = ""
dolla [x, y, z, a, b, c, d, e, f] = wrd [x, y, z] ++ "Million" ++ dolla [a, b, c, d, e, f]
dolla [x, y, z, a, b, c]          = wrd [x, y, z] ++ "Thousand" ++ dolla [a, b, c]
dolla [x, y, z]                   = wrd [x, y, z]

splitDig  :: Int -> [Int]
splitDig x = [div x 100000000, mod (div x 10000000) 10, mod (div x 1000000) 10, mod (div x 100000) 10, mod (div x 10000) 10, mod (div x 1000) 10, mod (div x 100) 10, mod (div x 10) 10, mod x 10]

dollar  :: Int -> String
dollar x | x == 0    = "ZeroDollars"
         | otherwise = dolla (splitDig x) ++ "Dollars"

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map (dollar . read) $ lines input
