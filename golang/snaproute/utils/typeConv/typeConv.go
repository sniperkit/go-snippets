//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package typeConv

import (
	"errors"
	"fmt"
)

type Error string

func ConvertToStrings(intf interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	switch intf := intf.(type) {
	case []interface{}:
		result := make([]string, len(intf))
		for i := range intf {
			if intf[i] == nil {
				continue
			}
			p, ok := intf[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("unexpected element type for Strings, got type %T", intf[i])
			}
			result[i] = string(p)
		}
		return result, nil
	case nil:
		return nil, errors.New("nil returned")
	case Error:
		return nil, errors.New("nil returned")
	}
	return nil, fmt.Errorf("unexpected type for Strings, got type %T", intf)
}

func ConvertToString(intf interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch intf := intf.(type) {
	case []byte:
		return string(intf), nil
	case string:
		return intf, nil
	case nil:
		return "", errors.New("nil returned")
	case Error:
		return "", errors.New("nil returned")
	}
	return "", fmt.Errorf("unexpected type for String, got type %T", intf)
}
