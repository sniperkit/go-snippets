import Data.Binary (decode)
import Data.Binary.Put (putWord16host, runPut)
import Data.Word (Word8)

littleEndian :: Bool
littleEndian = (decode . runPut $ putWord16host 42 :: Word8) == 42

main :: IO ()
main | littleEndian = putStrLn "LittleEndian"
     | otherwise    = putStrLn "BigEndian"
