package main

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
)

type Activity struct {
	Id         int64           `json:"id"`
	TemplateId string          `json:"template_id"`
	Event      ActivityPayload `json:"event"`
	TimeStamp
}

type ActivityList []Activity

func (a *ActivityList) List() error {
	_, err := dbmap.Select(a, "SELECT * FROM ACTIVITIES ORDER BY CreatedAt DESC")
	return err
}

func (a *Activity) Create() error {
	// run the UpdateTime ethod on the user model
	a.UpdateTime()

	// run the DB insert function
	err := dbmap.Insert(a)
	return err
}

func ActivityParser(p []ActivityPayload) []Activity {
	activities := []Activity{}

	for _, payload := range p {
		for _, template := range templates {
			activity := Activity{}

			if CheckEventType(template.Event, payload) && ValidateConditional(template.Condition, payload) {
				fmt.Println("Building activity")

				//Create the activity for this template and append to list
				activity.TemplateId = template.Id
				activity.Event = payload
				activity.UpdateTime()

				err := activity.Create()
				checkErr(err, "Problem saving activity")

				activities = append(activities, activity)

			}
		}
	}

	return activities
}

func CheckEventType(e string, p ActivityPayload) bool {
	if e == "*" {
		return true
	}
	return e == p.EventType
}

func ValidateConditional(c string, p ActivityPayload) bool {
	if c == "" {
		return true
	}

	buf := new(bytes.Buffer)

	tmpl, err := template.New("condition").Parse(fmt.Sprintf("{{%s}}", c))
	checkErr(err, "Failed to parse template condition")

	err = tmpl.Execute(buf, p)
	checkErr(err, "Failed to execute template condition")

	resultStr := buf.String()
	result, err := strconv.ParseBool(resultStr)
	checkErr(err, "Failed to parse bool from condition")

	return result
}
