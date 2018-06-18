package main

import (
    "fmt"
    "encoding/json"
    "log"
)

func main() {
   type Color struct {
       Space string
       Point json.RawMessage
   }

   type RGB struct {
       R uint8
       G uint8
       B uint8
   }

   type YCbCr struct {
       Y uint8
       Cb int8
       Cr int8
   }

   var j = []byte(`[
       {"Space": "RGB", "Point": {"R": 98, "G": 90, "B": 100}},
       {"Space": "YCbCr", "Point": {"Y": 255, "G": 0, "B": -10}}
   ]`)    

   var colors []Color
   err := json.Unmarshal(j, &colors)
   if err != nil {
       log.Fatalln("error ", err)
   }

   for _, color := range colors {
       var dst interface{}
       switch color.Space {
       case "RGB":
           dst = new(RGB)
       case "YCbCr":
           dst = new(YCbCr)
       }
       err := json.Unmarshal(color.Point, dst)
       if err != nil {
           log.Fatalln("error: ", err)
       }
       fmt.Printf("space is %s, point is %v \n", color.Space, dst)
   }
}
