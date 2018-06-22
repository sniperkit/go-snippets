package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	Id        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Owner     string        `json:"owner" bson:"owner"`
	Members   []string      `json:"members" bson:"members"`
	Meetings  Calendar      `json:"meetings" bson:"meetings"`
	Issues    Calendar      `json:"issues" bson:"issues"`
	Timesheet Calendar      `json:"timesheet" bson:"timesheet"`
	Invoices  Calendar      `json:"invoices" bson:"invoices"`
}

func (p *Project) Get(projectName, username string) error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Find(bson.M{"name": projectName}).One(&p)
	if err != nil {
		return err
	}

	if p.Owner != username {
		membership := false

		// Search for membership
		for _, m := range p.Members {
			if m == username {
				membership = true
			}
		}

		if !membership {
			p = nil
		}
	}

	return nil
}

func GetByUsername(username string) ([]Project, error) {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	var projects []Project
	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Find(nil).All(&projects)
	if err != nil {
		return nil, err
	}

	var myProjects []Project
	for _, p := range projects {
		if p.Owner == username {
			myProjects = append(myProjects, p)
		} else {
			// Search for membership
			for _, m := range p.Members {
				if m == username {
					myProjects = append(myProjects, p)
					break
				}
			}
		}
	}

	return myProjects, nil
}

func (p *Project) UpdateMeetings() error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Update(bson.M{"_id": p.Id}, bson.M{"$set": bson.M{"meetings": p.Meetings}})
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) UpdateIssues() error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Update(bson.M{"_id": p.Id}, bson.M{"$set": bson.M{"issues": p.Issues}})
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) UpdateTimesheet() error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Update(bson.M{"_id": p.Id}, bson.M{"$set": bson.M{"timesheet": p.Timesheet}})
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) UpdateInvoices() error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C("projects")
	err = c.Update(bson.M{"_id": p.Id}, bson.M{"$set": bson.M{"invoices": p.Invoices}})
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) Insert() error {
	return Insert("projects", interface{}(p))
}
