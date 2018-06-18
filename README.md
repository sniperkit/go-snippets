## Run

Run with:

```
./buildAndRun.sh
```

To activate the plugins, move the `.so` files from the
`plugins.available` directory to the `plugins` directory.

## Todo

- Malicious plugin -> read a folder and write to image
  - Check source code?
- Panic in plugin?
- In presentation
  - Add background about the solutions before the plugins: see the
    presentation from the dotGo.

## Presentation

Online here: https://gitpitch.com/gildasch/go-plugins-demo

- Add Go code at runtime
- Works only on linux for now
- Based on `dlfcn.h` (dynamic linking) with `cgo`
  - See [`src/plugin/plugin_dlopen.go`](https://tip.golang.org/src/plugin/plugin_dlopen.go)

- Only two functions: `Open(path string) (*Plugin, error)` and `(p
  *Plugin) Lookup(symName string) (Symbol, error)`
- Two types: `Plugin`, `Symbol`
- Nothing special on the code of the plugin, only the
  `-buildmode=plugin` build option

Example (from [the doc](https://tip.golang.org/pkg/plugin/)):
```
p, err := plugin.Open("plugin_name.so")
if err != nil {
	panic(err)
}
v, err := p.Lookup("V")
if err != nil {
	panic(err)
}
f, err := p.Lookup("F")
if err != nil {
	panic(err)
}
*v.(*int) = 7
f.(func())() // prints "Hello, number 7"
```
