package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"robostats/models/beat"
	"robostats/models/class"
	"robostats/models/instance"
	"robostats/models/session"
	"robostats/models/user"
	"strings"
	"time"
	"upper.io/i/v1/session/tokener"
)

const (
	randomInstancesLen = 5
	randomSessionsLen  = 10
	randomBeatsLen     = 1000
)

var (
	startTimeSpan int
	endTimeSpan   int
)

func init() {
	rand.Seed(time.Now().UnixNano())

	startTimeSpan = 30 * 24 * 3600 * 1000
	endTimeSpan = 15 * 24 * 3600 * 1000
}

func main() {

	// Removing data.
	user.UserCollection.Truncate()
	class.ClassCollection.Truncate()
	instance.InstanceCollection.Truncate()
	session.SessionCollection.Truncate()
	beat.BeatCollection.Truncate()

	// Adding user.
	u := &user.User{
		Email:    "user@example.com",
		Password: "pass",
	}

	if err := u.Create(); err != nil {
		log.Fatal(err)
	}

	// Adding some Classes
	classes := []string{
		"AR Parrot",
		"Toy Tank",
		"Raspberry PI",
		"Home computers",
		"Lab computers",
		"My office",
		"Handheld devices",
	}

	classesID := map[string]bson.ObjectId{}

	for _, className := range classes {
		k := class.Class{
			Name:   className,
			UserID: u.ID,
		}
		if err := k.Create(); err != nil {
			log.Fatal(err)
		}
		classesID[className] = k.ID
	}

	classesLen := len(classes)

	// Adding Instances to random classes.
	instances := make([]*instance.Instance, 0, randomInstancesLen)

	for i := 0; i < randomInstancesLen; i++ {
		// Pick a random class
		randomClassID := classesID[classes[rand.Intn(classesLen)]]
		// Create a random Intance.
		newInstance := &instance.Instance{
			UserID:  u.ID,
			ClassID: randomClassID,
			Data: map[string]string{
				"serial_number": strings.ToUpper(tokener.String(8)),
			},
		}
		if err := newInstance.Create(); err != nil {
			log.Fatal(err)
		}
		instances = append(instances, newInstance)
	}

	// Adding random sessions associated to classes.
	sessions := make([]*session.Session, 0, randomSessionsLen)
	for i := 0; i < randomSessionsLen; i++ {
		randomInstance := instances[rand.Intn(randomInstancesLen)]

		randomStartTime := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local)

		randomStartTime = randomStartTime.Add(time.Millisecond * time.Duration(rand.Intn(startTimeSpan)))
		randomEndTime := randomStartTime.Add(time.Millisecond * time.Duration(rand.Intn(endTimeSpan)))

		newSession := &session.Session{
			UserID:     u.ID,
			ClassID:    randomInstance.ClassID,
			InstanceID: randomInstance.ID,
			SessionKey: tokener.String(20),
			StartTime:  randomStartTime,
			EndTime:    randomEndTime,
		}

		if err := newSession.Create(); err != nil {
			log.Fatal(err)
		}
		sessions = append(sessions, newSession)
	}

	// Adding randomly generated beats to instances.
	beatTimes := map[bson.ObjectId]int{}
	beatLatLng := map[bson.ObjectId][2]float64{}

	for i := 0; i < randomBeatsLen; i++ {
		randomSession := sessions[rand.Intn(randomSessionsLen)]

		if _, ok := beatLatLng[randomSession.ID]; !ok {
			beatLatLng[randomSession.ID] = [2]float64{
				19.42705,
				-99.12757,
			}
		}

		beatLatLng[randomSession.ID] = [2]float64{
			beatLatLng[randomSession.ID][0] + float64(rand.Intn(3)-1)*float64(rand.Intn(99999))/100000,
			beatLatLng[randomSession.ID][1] + float64(rand.Intn(3)-1)*float64(rand.Intn(99999))/100000,
		}

		if _, ok := beatTimes[randomSession.ID]; !ok {
			beatTimes[randomSession.ID] = 0
		}

		beatTimes[randomSession.ID] += 1 + rand.Intn(5)

		newBeat := &beat.Beat{
			UserID:     u.ID,
			ClassID:    randomSession.ClassID,
			InstanceID: randomSession.InstanceID,
			SessionID:  randomSession.ID,
			Data: bson.M{
				"cpu":    float64(0.8 + rand.Float32()*0.3),
				"height": float64(rand.Intn(3)) + rand.Float64(),
			},
			LocalTime: beatTimes[randomSession.ID],
			LatLng:    beatLatLng[randomSession.ID],
		}

		if err := newBeat.Create(); err != nil {
			log.Fatal(err)
		}

	}

}
