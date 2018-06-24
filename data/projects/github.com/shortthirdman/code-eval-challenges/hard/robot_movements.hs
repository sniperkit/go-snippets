import Data.Bits (setBit, testBit)

robot      :: Int -> Int -> Int -> Int
robot _ 3 3 = 1
robot f x y = n + e + s + w
            where n | y > 0 && not (testBit f (x+4*(y-1))) = robot (setBit f (x+4*(y-1))) x (y-1)
                    | otherwise                            = 0
                  e | x < 3 && not (testBit f (x+1+4*y))   = robot (setBit f (x+1+4*y)) (x+1) y
                    | otherwise                            = 0
                  s | y < 3 && not (testBit f (x+4*(y+1))) = robot (setBit f (x+4*(y+1))) x (y+1)
                    | otherwise                            = 0
                  w | x > 0 && not (testBit f (x-1+4*y))   = robot (setBit f (x-1+4*y)) (x-1) y
                    | otherwise                            = 0

main :: IO ()
main = print $ robot 1 0 0
