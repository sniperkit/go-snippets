// START OMIT
/***********************
 * Mission: All Aboard *
 ***********************
 *
 * Logging into the agency network, you're quickly into the operation manifest code.
 * Sure enough, you're listed with no equipment, a typical agency snafu.
 * You're going to need full tactical gear for this mission.
 * Let's go ahead and fix that.
 *
 */

package main

import(
  "fmt"
  "strconv"
)

func main() {
  total_agents := 5
  agents := make([]Agent,0)
  agents = append(agents, Agent{name: "Jason Marshall", equipment: "none"})
  agents = append(agents, Agent{name: "Stone Boswell", equipment: "full"})
  agents = append(agents, Agent{name: "Jane Johnson", equipment: "full"})
  agents = append(agents, Agent{name: "Max Carter", equipment: "full"})
  agents = append(agents, Agent{name: "Kay White", equipment: "full"})
// EDITABLE OMIT

  // Your code

// UNEDITABLE OMIT
  gear_table := createGearTable(agents, total_agents)
  fmt.Println(gear_table)
}

func createGearTable(agents []Agent, total_agents int) string {
  if len(agents) != total_agents {
    return "ERROR: Wrong number of agents"
  }
  gear_table := "Operation Go: Agent Manifest\n"
  gear_table += "----------------------------"
  for num, agent := range agents {
    gear_table += "\n" + strconv.Itoa(num+1) + ". " + agent.name
    gear_table += ", Gear: " + agent.equipment
  }
  return gear_table
}

type Agent struct {
  name string
  equipment string
}
// END OMIT
