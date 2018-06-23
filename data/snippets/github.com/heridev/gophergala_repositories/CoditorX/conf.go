// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/b3log/wide/log"
	"github.com/b3log/wide/util"
)

// Logger
var logger *log.Logger

// Configuration.
type config struct {
	IP                    string // server ip, ${ip}
	Port                  string // server port
	Context               string // server context
	Server                string // server host and port ({IP}:{Port})
	StaticServer          string // static resources server scheme, host and port (http://{IP}:{Port})
	LogLevel              string // logging level: trace/debug/info/warn/error
	Channel               string // channel (ws://{IP}:{Port})
	HTTPSessionMaxAge     int    // HTTP session max age (in seciond)
	StaticResourceVersion string // version of static resources
	MaxProcs              int    // Go max procs
	RuntimeMode           string // runtime mode (dev/prod)
	WD                    string // current working direcitory, ${pwd}
	Locale                string // default locale
	Workspace             string // the space to store the user's file
}

var conf *config

func loadConf(confIP, confPort, confChannel string) {
	bytes, err := ioutil.ReadFile("conf/coditor.json")
	if nil != err {
		logger.Error(err)

		os.Exit(-1)
	}

	conf = &config{}

	err = json.Unmarshal(bytes, conf)
	if err != nil {
		logger.Error("Parses [coditor.json] error: ", err)

		os.Exit(-1)
	}

	log.SetLevel(conf.LogLevel)

	logger = log.NewLogger(os.Stdout)

	logger.Debug("Conf: \n" + string(bytes))

	// Working Driectory
	conf.WD = util.OS.Pwd()
	logger.Debugf("${pwd} [%s]", conf.WD)

	// IP
	ip, err := util.Net.LocalIP()
	if err != nil {
		logger.Error(err)

		os.Exit(-1)
	}

	logger.Debugf("${ip} [%s]", ip)

	conf.IP = strings.Replace(conf.IP, "${ip}", ip, 1)
	if "" != confIP {
		conf.IP = confIP
	}

	if "" != confPort {
		conf.Port = confPort
	}

	// Server
	conf.Server = strings.Replace(conf.Server, "{IP}", conf.IP, 1)

	// Static Server
	conf.StaticServer = strings.Replace(conf.StaticServer, "{IP}", conf.IP, 1)
	conf.StaticServer = strings.Replace(conf.StaticServer, "{Port}", conf.Port, 1)

	conf.StaticResourceVersion = strings.Replace(conf.StaticResourceVersion, "${time}", strconv.FormatInt(time.Now().UnixNano(), 10), 1)

	// Channel
	conf.Channel = strings.Replace(conf.Channel, "{IP}", conf.IP, 1)
	conf.Channel = strings.Replace(conf.Channel, "{Port}", conf.Port, 1)

	if "" != confChannel {
		conf.Channel = confChannel
	}

	conf.Server = strings.Replace(conf.Server, "{Port}", conf.Port, 1)

}
