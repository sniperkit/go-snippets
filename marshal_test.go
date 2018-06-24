package xmp

import (
	"fmt"
	"testing"
	"time"
)

func TestAsPOD(t *testing.T) {
	created, _ := time.Parse(time.RFC3339, "2016-04-03T14:29:11-04:00")
	expect := &POD{
		Title:       "What Climate Change Means for Massachusetts",
		Description: "This fact sheet provides a concise overview of the observed and projected effects and impacts of climate change on Massachusetts.",
		Created:     &created,
		Keyword:     []string{"EPA", "climate change", "state", "impacts", "fact sheet", "summary"},
		Publisher:   "US EPA; OAR; Climate Change Division",
	}
	packet, err := Unmarshal([]byte(xmpExample))
	if err != nil {
		t.Error(err.Error())
		return
	}

	pod := packet.AsPOD()

	if err := ComparePod(expect, pod); err != nil {
		t.Errorf(err.Error())
	}
}

func ComparePod(a, b *POD) error {
	if a == nil && b != nil || a != nil && b == nil {
		return fmt.Errorf("nil mismatch, %s != %s", a, b)
	}

	if a.Title != b.Title {
		return fmt.Errorf("Title mismatch: %s != %s", a.Title, b.Title)
	}

	if a.Description != b.Description {
		return fmt.Errorf("Description mismatch: %s != %s", a.Description, b.Description)
	}

	if a.Publisher != b.Publisher {
		return fmt.Errorf("Publisher mismatch: %s != %s", a.Publisher, b.Publisher)
	}

	if err := CompareStringSlices(a.Keyword, b.Keyword); err != nil {
		return fmt.Errorf("keyword mismatch: %s", err.Error())
	}

	return nil
}

func CompareTimePointers(a, b *time.Time) error {
	if a == nil && b != nil || a != nil && b == nil {
		return fmt.Errorf("nil mismatch: %s != %s", a, b)
	}

	if !a.Equal(*b) {
		return fmt.Errorf("time mismatch: %s != %s", a, b)
	}
	return nil
}

func CompareStringSlices(a, b []string) error {
	if a != nil && b == nil || a == nil && b != nil {
		return fmt.Errorf("nil mismatch: %s != %s", a, b)
	}
	if a != nil && b != nil {
		if len(a) != len(b) {
			return fmt.Errorf("slice length mismatch: %d != %d", len(a), len(b))
		}

		for i, ai := range a {
			if ai != b[i] {
				return fmt.Errorf("slice index %d entry mismatch: %s != %s", i, ai, b[i])
			}
		}
	}
	return nil
}
