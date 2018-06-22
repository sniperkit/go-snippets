// START OMIT
/***********************
 * Mission: Where Am I *
 ***********************
 *
 *
 *
 */

package main

import(
  "fmt"
  "strings"
)

func main() {

  request := Request{agent: "Jason Marshall", mission: "Operation Go"}

// EDITABLE OMIT

  // Your code


// UNEDITABLE OMIT

  fmt.Println("Location:", pingSattelite(request))
}

type Request struct {
  agent string
  mission string
}

func pingSattelite(request Request) string {
  if accessible(request.agent) {
    return revealLocation()
  }
  return "ACCESS DENIED"
}

func accessible(agent string) bool {
  a := strings.Split(agent, " ")
  if len(a) != 3 { return false }
  b := a[0]
  c := a[1]
  d := a[2]
  x := strings.EqualFold(b, c)
  y := b != strings.ToLower(c)
  z := strings.Index(d, b+c) == 1 && len(d) == 5
  return x && y && z
}

// END OMIT
func revealLocation() string {
  return "16.7333,-169.5274"
}
