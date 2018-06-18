package generic

import (
	"github.com/go-zero-boilerplate/facades/mitigation"
	"github.com/go-zero-boilerplate/facades/tasks"
)

//TryToMitigateAndAlert will if the mitigateCondition is TRUE call the mitigateFunc and log the answer getMitigationMessage()
func TryToMitigateAndAlert(mitigateCondition bool, mitigateFunc func() error, getMitigationMessage func() string) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		if err := mitigation.TryAndAlert(mitigateCondition, mitigateFunc, getMitigationMessage); err != nil {
			taskCtx.Error = tasks.NewErrorDefaultStatus(err)
			return
		}
	}
}
