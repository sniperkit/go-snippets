// START OMIT
/*************************
 * Mission: Time to Jump *
 *************************
 *
 * That's weird, looking at this drop sequence, it looks like
 * you would miss the island...
 * You know that you're cruising at 3000m and that the island is 1000m
 * directly northwest from the jump site, with winds blowing with you
 * You decide to fix the code to ensure you'll land safely.
 *
 */

package main

import(
  "fmt"
)

func main() {

  elevation := 4000
  destination_north := 1000
  destination_west := 1000

  wind_north := -4
  wind_west := -4

// EDITABLE OMIT

  // Your code

  deploy_chute_at_elevation := 1000

// UNEDITABLE OMIT
  landing_location := simDrop(elevation,
                          destination_north,
                          destination_west,
                          wind_north,
                          wind_west,
                          deploy_chute_at_elevation)
  fmt.Println("Initiate deploy at:", deploy_chute_at_elevation, "meters -", landing_location)
}

func simDrop(elevation int,
             dest_north int,
             dest_west int,
             wind_north int,
             wind_west int,
             deploy_chute_at_elevation int) string {

  location_north := 0
  location_west := 0

  deployed := false
  for elevation > 0 {
    if elevation < deploy_chute_at_elevation {
      deployed = true
    }
    if !deployed {
      location_north -= wind_north / 4
      location_west -= wind_west / 4
    } else {
      location_north -= wind_north
      location_west -= wind_west
    }
    elevation -= 10
  }
  landing_location := "You missed the target!"
  if location_north == dest_north && location_west == dest_west {
    landing_location = "You hit the target!"
  }

  return landing_location
}

// END OMIT
