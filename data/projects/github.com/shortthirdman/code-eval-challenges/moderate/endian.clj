(println
  (if (= (System/getProperty "sun.cpu.endian") "little")
    "LittleEndian" "BigEndian"))
