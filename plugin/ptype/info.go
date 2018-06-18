package ptype

type PluginInfo struct {
	Uri     string
	Version string
	Descrip string
	Name    string
	Command string
	Usage   string
}

type PluginDeclarer interface {
	DeclarePlugin() PluginInfo
}
