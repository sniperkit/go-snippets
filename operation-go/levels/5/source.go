// START OMIT
/**************************
 * Mission: Into the Open *
 **************************
 *
 * The passcode looks challenging, but with enough skill you
 * might be able to break through.  It looks like the codes are in
 * the format of XX-XX-XX and after three attempts a brute force blocker
 * kicks in.
 *
 */

package main

import(
  "fmt"
  "strconv"
)

func main() {

  // EDITABLE OMIT

  guesses := make([]string,0)
  guesses = append(guesses, "00-00-00")
  guesses = append(guesses, "01-00-00")
  guesses = append(guesses, "02-00-00")

  // UNEDITABLE OMIT
  response := validateCode(guesses)
  fmt.Println("Code:", response)
}

func validateCode(codes []string) string {
  for _, c := range codes {
    if isCorrect(c) { return c }
  }
  // Denies more than three attempts
  return foilBruteForce()
}

// END OMIT

func isCorrect(guess string) bool {
  code := "11-10-09"
  if guess == code {
    return true
  }
  return false
}

var attempts int = 0

func foilBruteForce() string {
  attempts++
  if attempts > 3 {
    return "LOCKED: Tried more than three times"
  }
  return "Wrong Guess" + strconv.Itoa(attempts)
}
