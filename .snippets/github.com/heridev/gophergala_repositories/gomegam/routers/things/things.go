/*
** Copyright [2012-2014] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package things

import (
	"github.com/gophergala/gomegam/routers/base"
	"strings"
	"fmt"
	"encoding/json"
	"github.com/gophergala/gomegam/app"
	"github.com/gophergala/gomegam/global"
)

type ThingsRouter struct {
	base.BaseRouter
}

func (this *ThingsRouter) Get() {

    result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

    req := this.Ctx.Input.Param(":id")
    data := strings.Split(req, "-")
    
    d := global.App{DeviceID: data[0], Hash: data[1]}
    b, err := json.Marshal(d)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(b))
    
    err = app.AnalyticsProcess(&d)
		if err != nil {
			fmt.Println(err)
		}
  	result["success"] = true
	result["data"] = string(b)
}


