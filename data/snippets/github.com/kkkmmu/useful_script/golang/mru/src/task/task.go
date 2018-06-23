package task

import (
	"condition"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"logger"
	"net/url"
	"routine"
	"rut"
	"taskresult"
)

type Task struct {
	Name          string               `json:"name"`
	PreCondition  *condition.Condition `json:"precondition"`
	Routine       *routine.Routine     `json:"routine"`
	PostCondition *condition.Condition `json:"postcondition"`
	Clear         *routine.Routine     `json:"clear"`
	Description   string               `json:"description"`
	ID            string               `json:"id"`
}

func (t Task) String() string {
	res := fmt.Sprintf("Task: %20s: \n", t.Name)
	res += fmt.Sprintf("         %#v", t.PreCondition)
	res += fmt.Sprintf("         %#v", t.Routine)
	res += fmt.Sprintf("         %#v", t.PreCondition)
	res += fmt.Sprintf("         %#v", t.Clear)
	res += fmt.Sprintf("         %#v", t.Description)
	res += fmt.Sprintf("         %s", t.ID)
	return res
}

func Hash(name []byte) []byte {
	hash := sha1.New()
	return []byte(hex.EncodeToString(hash.Sum([]byte("taskTASK" + string(name)))))
}

func (t *Task) Run(ctx context.Context, db *rut.DB) *taskresult.Result {
	//This must be at first, I want to do clear work when error happend.
	defer t.RunClearRoutine(ctx, db)
	message := fmt.Sprintf("[Running Task]: {%s}\n", t.Name)
	logger.Push(ctx, message)

	var res *taskresult.Result

	if res = t.CheckPreCondition(ctx, db); !res.Success {
		logger.Push(ctx, res.String())
		return res
	}

	if res = t.RunMainRoutine(ctx, db); !res.Success {
		logger.Push(ctx, res.String())
		return res
	}

	if res = t.CheckPostCondition(ctx, db); !res.Success {
		logger.Push(ctx, res.String())
		return res
	}

	res = &taskresult.Result{Name: t.Name, Description: t.Description, Success: true}
	logger.Push(ctx, res.String())
	return res
}

func (t *Task) CheckPreCondition(ctx context.Context, db *rut.DB) *taskresult.Result {
	for _, r := range db.DB {
		r.GoInitMode()
	}

	if err := t.PreCondition.Check(ctx, db); err != nil {
		return &taskresult.Result{
			Name:        t.Name,
			Description: t.Description,
			Success:     false,
			Message:     err.Error(),
		}
	}

	return &taskresult.Result{
		Name:        t.Name,
		Description: t.Description,
		Success:     true,
	}
}

func (t *Task) CheckPostCondition(ctx context.Context, db *rut.DB) *taskresult.Result {
	for _, r := range db.DB {
		r.GoInitMode()
	}

	if err := t.PostCondition.Check(ctx, db); err != nil {
		return &taskresult.Result{
			Name:        t.Name,
			Description: t.Description,
			Success:     false,
			Message:     err.Error(),
		}
	}

	return &taskresult.Result{
		Name:        t.Name,
		Description: t.Description,
		Success:     true,
	}

}

func (t *Task) RunMainRoutine(ctx context.Context, db *rut.DB) *taskresult.Result {
	for _, r := range db.DB {
		r.GoInitMode()
	}

	if err := t.Routine.Run(ctx, db); err != nil {
		return &taskresult.Result{
			Name:        t.Name,
			Description: t.Description,
			Success:     false,
			Message:     err.Error(),
		}
	}

	return &taskresult.Result{
		Name:        t.Name,
		Description: t.Description,
		Success:     true,
	}

}

func (t *Task) RunClearRoutine(ctx context.Context, db *rut.DB) *taskresult.Result {
	for _, r := range db.DB {
		r.GoInitMode()
	}
	if err := t.Clear.Run(ctx, db); err != nil {
		logger.Push(ctx, err.Error())
		return &taskresult.Result{
			Name:        t.Name,
			Description: t.Description,
			Success:     false,
			Message:     err.Error(),
		}
	}

	return &taskresult.Result{
		Name:        t.Name,
		Description: t.Description,
		Success:     true,
	}

}

func (t *Task) SetPreCondition(con *condition.Condition) {
	t.PreCondition = con
}

func (t *Task) SetPostCondition(con *condition.Condition) {
	t.PostCondition = con
}

func (t *Task) SetMainRoutine(r *routine.Routine) {
	t.Routine = r
}

func (t *Task) SetClearRoutine(r *routine.Routine) {
	t.Clear = r
}

func (t *Task) SetName(name string) {
	t.Name = name
}

func (t *Task) SetDescription(desc string) {
	t.Description = desc
}

func IsTaskParamsValid(in url.Values) bool {
	for k, v := range in {
		log.Println(k, "----------->", v, len(v))
		if len(v) == 0 {
			return false
		}
	}
	return true
}
