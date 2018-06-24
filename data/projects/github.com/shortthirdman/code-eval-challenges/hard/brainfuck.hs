import System.Environment (getArgs)
import Data.Char (chr)
import Data.Map (Map)
import qualified Data.Map as Map

bmap        :: Int -> [Int] -> String -> [(Int, Int)]
bmap x ys zs | x == length zs = []
             | c == '['       = bmap (succ x) (x:ys) zs
             | c == ']'       = (x, head ys + 1) : (head ys, x + 1) : bmap (succ x) (tail ys) zs
             | otherwise      = bmap (succ x) ys zs
             where c = zs!!x

bracketmap   :: String -> Map Int Int
bracketmap xs = Map.fromList (bmap 0 [] xs)

bfuck                  :: Int -> Int -> Map Int Int -> Map Int Int -> String -> String
bfuck x y brack dats zs | x == length zs      = ""
                        | c == '+' && d < 255 = bfuck (succ x) y brack (Map.insert y (succ d) dats) zs
                        | c == '+'            = bfuck (succ x) y brack (Map.insert y 0 dats) zs
                        | c == '-' && d > 0   = bfuck (succ x) y brack (Map.insert y (pred d) dats) zs
                        | c == '-'            = bfuck (succ x) y brack (Map.insert y 255 dats) zs
                        | c == '>'            = bfuck (succ x) (succ y) brack dats zs
                        | c == '<'            = bfuck (succ x) (pred y) brack dats zs
                        | c == '.'            = chr d : bfuck (succ x) y brack dats zs
                        | c == ','            = bfuck (succ x) y brack dats zs
                        | c == '[' && d == 0  = bfuck (Map.findWithDefault 0 x brack) y brack dats zs
                        | c == '['            = bfuck (succ x) y brack dats zs
                        | c == ']' && d /= 0  = bfuck (Map.findWithDefault 0 x brack) y brack dats zs
                        | otherwise           = bfuck (succ x) y brack dats zs
                        where c = zs!!x
                              d = Map.findWithDefault 0 y dats

brainfuck   :: String -> String
brainfuck xs = bfuck 0 0 (bracketmap xs) Map.empty xs

main :: IO ()
main = do
    [inpFile] <- getArgs
    input <- readFile inpFile
    putStr . unlines . map brainfuck $ lines input
