package util

import "fmt"

func Await(cres chan interface{}, cerr chan error) (result interface{}, err error) {
	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("Await Failed: channels are undefined")
	}

	resolved := false
	for !resolved {
		select {
		case res, ok := <-cres:
			if ok {
				result = res
				resolved = true
			}
		case error, ok := <-cerr:
			if ok {
				err = error
				resolved = true
			}
		}
	}

	return result, err
}

func AwaitError(cerr chan error) (err error) {
	if cerr == nil {
		return fmt.Errorf("Await Failed: channel is undefined")
	}

	err, _ = <-cerr

	return err
}
