package main

import (
      "github.com/tarm/goserial"
      "log"
      "net/http"
)

func main() {
  c1 := new(serial.Config)
  c1.Name = "/dev/tty.usbmodem1421"
  c1.Baud = 9600
  
  s, err := serial.OpenPort(c1)
  if err != nil {
    log.Fatal(err)
  }

  n := 0
  buf := make([]byte, 256)
  for {
    n, err = s.Read(buf)
    if err != nil {
      log.Fatal(err)
       s.Close()
       break;
    }
    if n>10 {
      resp, _ := http.Get("http://findaparking.appspot.com/parkAuto")
    }
    time.Sleep(3000 * time.Millisecond)
  }
  
      

 
}