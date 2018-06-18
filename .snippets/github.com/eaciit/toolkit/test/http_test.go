package toolkittest

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/eaciit/toolkit"
)

func Test_Http(t *testing.T) {
	tstart := time.Now()
	url := "http://localhost"

	type Obj struct {
		Wg      *sync.WaitGroup
		Success int
		Errors  []string
	}
	goobj := &Obj{new(sync.WaitGroup), 0, []string{}}
	for i := 0; i < 9500; i++ {
		goobj.Wg.Add(1)
		go func(goobj *Obj, i int) {
			//func(goobj *Obj, i int) {
			defer goobj.Wg.Done()
			r, e := HttpCall(url, "GET", nil, nil)
			if e != nil {
				goobj.Errors = append(goobj.Errors,
					fmt.Sprintf("(%d) Unable to call %s, got %s", i, url, e.Error()))
				//t.Errorf("(%d) Unable to call %s, got %s", i, url, e.Error())
				//return
			} else {
				s := HttpContentString(r)
				if s == "" {
					goobj.Errors = append(goobj.Errors,
						fmt.Sprintf("(%d) Got nothing", i))
				} else {
					goobj.Success++
				}
			}
		}(goobj, i)
	}
	goobj.Wg.Wait()

	if len(goobj.Errors) == 0 {
		fmt.Printf("Done in %v \n", time.Since(tstart))
	} else {
		for _, e := range goobj.Errors {
			fmt.Println(e)
		}
		t.Errorf("Fail due to above errors. Success %d Fails %d", goobj.Success, len(goobj.Errors))
	}

}
