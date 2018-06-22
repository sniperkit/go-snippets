package client

import (
	"log"
	"testing"
)

func TestClientLogin(t *testing.T) {
	var err error

	client := Client{
		Email:    "user@example.com",
		Password: "pass",
	}

	if err = client.Login(); err != nil {
		t.Fatal(err)
	}

	var classes []Class
	if classes, err = client.GetClasses(); err != nil {
		t.Fatal(err)
	}

	for _, class := range classes {
		// Add new instance
		var i *Instance
		if i, err = client.RegisterInstance(class.ID, map[string]string{"foo": "bar"}); err != nil {
			t.Fatal(err)
		}

		// Add new session to new instance
		var s *Session
		if s, err = client.RegisterSession(i.ID, map[string]string{"nuff": "said"}); err != nil {
			t.Fatal(err)
		}

		// Push an event.
		var e *Event
		if e, err = s.PushEvent(map[string]float64{"cpu": 0.1227}); err != nil {
			t.Fatal(err)
		}

		log.Printf("Event added: %v\n", e)

		// List instances.
		var instances []Instance
		if instances, err = client.GetInstancesByClassID(class.ID); err != nil {
			t.Log(err)
		}
		for _, instance := range instances {
			var sessions []Session
			if sessions, err = client.GetSessionsByInstanceID(instance.ID); err != nil {
				t.Log(err)
			}
			log.Printf("Active sessions: %d\n", len(sessions))
		}
	}

}
