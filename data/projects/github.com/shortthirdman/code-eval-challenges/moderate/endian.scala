import java.nio.ByteOrder

object Main extends App {
  if (System.getProperty("sun.cpu.endian") == "little")
    println("LittleEndian")
  else
    println("BigEndian")
}
