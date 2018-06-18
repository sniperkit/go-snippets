package plugin
import p0 "github.com/gizak/gopp/plugin/gocolor"
import p1 "github.com/gizak/gopp/plugin/parsepanic"
import p2 "github.com/gizak/gopp/plugin/rm"
func init() {
 	Plugins = append(Plugins,PluginEntry{Name:"gocolor",Plugin:p0.PLUGIN,Uri:"github.com/gizak/gopp/plugin/gocolor"})
	Plugins = append(Plugins,PluginEntry{Name:"parsepanic",Plugin:p1.PLUGIN,Uri:"github.com/gizak/gopp/plugin/parsepanic"})
	Plugins = append(Plugins,PluginEntry{Name:"rmpkg",Plugin:p2.PLUGIN,Uri:"github.com/gizak/gopp/plugin/rm"})
}