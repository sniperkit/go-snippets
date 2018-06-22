/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	_ "expvar"
	"github.com/infinitbyte/framework"
	"github.com/infinitbyte/framework/core/module"
	pipe "github.com/infinitbyte/framework/core/pipeline"
	q "github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/framework/modules/api"
	"github.com/infinitbyte/framework/modules/filter"
	"github.com/infinitbyte/framework/modules/pipeline"
	"github.com/infinitbyte/framework/modules/queue"
	"github.com/infinitbyte/framework/modules/stats"
	"github.com/infinitbyte/framework/modules/storage"
	"github.com/infinitbyte/framework/modules/ui"
	"github.com/medcl/csv2sql/config"
	"github.com/medcl/csv2sql/plugin"
)

func main() {

	terminalHeader := ("  _____  _______      _____   _____  ____  _\n")
	terminalHeader += (" / ____|/ ____\\ \\    / /__ \\ / ____|/ __ \\| |\n")
	terminalHeader += ("| |    | (___  \\ \\  / /   ) | (___ | |  | | |\n")
	terminalHeader += ("| |     \\___ \\  \\ \\/ /   / / \\___ \\| |  | | |\n")
	terminalHeader += ("| |____ ____) |  \\  /   / /_ ____) | |__| | |____\n")
	terminalHeader += (" \\_____|_____/    \\/   |____|_____/ \\___\\_\\______|\n")

	terminalFooter := ("                         |    |                \n")
	terminalFooter += ("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) \n")
	terminalFooter += (" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| \n")
	terminalFooter += (" ____/                             ___/        \n")

	app := framework.NewApp("CSV2SQL", "An util to convert csv data to SQL scripts.",
		util.TrimSpaces(config.Version), util.TrimSpaces(config.LastCommitLog), util.TrimSpaces(config.BuildDate), terminalHeader, terminalFooter)

	app.Init(nil)
	defer app.Shutdown()

	app.Start(func() {
		//modules
		module.New()

		//load core modules first
		module.Register(module.Storage, storage.StorageModule{})
		module.Register(module.Filter, filter.FilterModule{})
		module.Register(module.Stats, stats.SimpleStatsModule{})
		module.Register(module.Queue, queue.DiskQueue{})
		module.Register(module.System, pipeline.PipelineFrameworkModule{})
		module.Register(module.API, api.APIModule{})
		module.Register(module.UI, ui.UIModule{})

		//load plugins
		module.RegisterPlugin(module.Tools, plugin.CSV2SQLPlugin{})

		//start modules
		module.Start()

		//manual trigger a pipeline event
		context := pipe.Context{}
		q.Push("primary", util.ToJSONBytes(context))

	})

}
