package main
import (
"fmt"
)
/*I thought this may be hard to follow...
So we figure the sheet area and figure each of the part areas.
Then we add parts till it fills the sheet as full as possible.*/
func calc(sWidth int8, sLength int8){
var sArea = sWidth*sLength
var pWidth int8
fmt.Println("Enter part width")
fmt.Scanf("%f", &pWidth)
var pLength int8
fmt.Println("Enter part length")
fmt.Scanf("%f", &pLength)
var pArea =pWidth*pLength
var aaa = 1
for p1:=pArea; p1<=sArea; p1++{
aaa++
}
fmt.Println(aaa)
}
func main() {
var length int8
fmt.Println("Enter sheet length")
fmt.Scanf("%f", &length)
var width int8
fmt.Println("Enter sheet width")
fmt.Scanf("%f", &width)
calc(width, length)
