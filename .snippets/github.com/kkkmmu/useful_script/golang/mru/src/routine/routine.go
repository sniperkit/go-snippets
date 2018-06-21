package routine

import (
	"assertion"
	"context"
	"errors"
	"fmt"
	"rut"
)

type Routine struct {
	Name        string                 `json:"name"`
	Assertions  []*assertion.Assertion `json:"assertions"`
	Description string                 `json:"description"`
}

func (r *Routine) Run(ctx context.Context, db *rut.DB) error {
	for _, a := range r.Assertions {
		msg, ok := a.Do(ctx, db)
		if !ok {
			return errors.New(fmt.Sprintf("Routine: %s failed with: %s", r.Name, msg))
		}
	}

	return nil
}
