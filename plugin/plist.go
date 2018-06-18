package plugin

import (
	"github.com/gizak/gopp/plugin/ptype"
	"log"
)

type PluginEntry struct {
	Name   string
	Plugin interface{}
	Uri    string
}

var Plugins = []PluginEntry{
	{Name: "plugin", Plugin: PLUGIN, Uri: "github.com/gizak/gopp/plugin"},
	{Name: "cmdhelp", Plugin: pCmdHelp, Uri: "github.com/gizak/gopp/plugin"}}

type PListEntry struct {
	Plugin string
	Uri    string
	Deps   []string
}

type PList []PListEntry

type Config struct {
	Plugins PList `yaml:"plugins"`
}

var OutputRewriters = []ptype.OutputRewriter{}

var SubcmdRunners = make(map[string]ptype.SubcmdRunner)

var PluginDecls = make(map[string]ptype.PluginInfo)

func regPlugins() {

	for _, p := range Plugins {
		if pp, ok := p.Plugin.(ptype.PluginDeclarer); ok {
			info := pp.DeclarePlugin()
			log.Println("REG PluginDeclarer:\t" + info.Name)
			PluginDecls[p.Name] = info
		}

		if pp, ok := p.Plugin.(ptype.SubcmdRunner); ok {
			log.Println("REG SubcmdRunner:  \t" + pp.Subcmd())
			SubcmdRunners[pp.Subcmd()] = pp
		}

		if pp, ok := p.Plugin.(ptype.OutputRewriter); ok {
			log.Println("REG OutputRewriter:\t" + p.Name)
			OutputRewriters = append(OutputRewriters, pp)
		}
	}
}
