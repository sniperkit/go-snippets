package condition

import (
	"assertion"
	"context"
	"errors"
	"fmt"
	"rut"
	"sync"
)

type Condition struct {
	Name        string                 `json:"name"`
	Assertions  []*assertion.Assertion `json:"assertions"`
	Description string                 `json:"description"`
}

func (cd Condition) String() string {
	return fmt.Sprintf("Condition: %15s, Assertions %v, Description: %s", cd.Name, cd.Assertions, cd.Description)
}

func (cd *Condition) Check(ctx context.Context, db *rut.DB) error {
	done := make(chan bool)
	err := make(chan string)

	defer close(done)
	defer close(err)

	go func(done chan<- bool, err chan<- string, as []*assertion.Assertion) {
		wg := sync.WaitGroup{}
		for _, a := range as {
			wg.Add(1)
			msg, ok := a.Do(ctx, db)
			if !ok {
				err <- fmt.Sprintf("Condition: %s check(%s) failed with: %s", cd.Name, cd.Description, msg)
				return
			}
			wg.Done()
		}
		wg.Wait()
		done <- true
	}(done, err, cd.Assertions)

	select {
	case <-done:
		return nil
	case msg := <-err:
		return errors.New(msg)
	}
}
