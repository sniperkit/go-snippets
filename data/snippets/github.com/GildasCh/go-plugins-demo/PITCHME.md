#HSLIDE

# Go plugins

Gildas Chabot, leboncoin

#HSLIDE

## Demo!

#HSLIDE

- Released in 1.8
- Add Go code at runtime
- Works only on linux for now
- Based on C's `dlfcn.h` (dynamic linking) with `cgo`
  - See [`src/plugin/plugin_dlopen.go`](https://tip.golang.org/src/plugin/plugin_dlopen.go)

#HSLIDE

## Need?

- Extend a service with third party functions
  - Web server
  - Media player
- Update behaviour at runtime

#HSLIDE

## Before

- Add plugins at compile time (Caddy)
- HTTP/RPC calls

- Check:
  - https://github.com/hashicorp/go-plugin
  - https://github.com/natefinch/pie


#HSLIDE

## Now: Plugin side

- Nothing special on the code of the plugin
  - `-buildmode=plugin` build option
- Creates a `.so`

#HSLIDE

## Now: Caller side

- Only two functions:
  - `Open(path string) (*Plugin, error)`
  - `(p  *Plugin) Lookup(symName string) (Symbol, error)`
- Two types:
  - `Plugin`
  - `Symbol`
- When the plugin if first `Open`, all _new_ packages have their
  `init` function called

See https://tip.golang.org/pkg/plugin/

#HSLIDE

## Example

```
p, err := plugin.Open("plugin.so")

fs, err := p.Lookup("Transform")
f, ok := fs.(func(image.Image) image.Image)

vs, err := p.Lookup("Priority")
v, ok = *vs.(*int)
```

#HSLIDE

## Safety?

- Plugin == safe?
- Same rights as caller <!-- .element: class="fragment" -->
- Check the sources of the plugins! <!-- .element: class="fragment" -->

#HSLIDE

## Stable?

- Invalid `.so` file
  - `fatal error: runtime: no plugin module data`
- Can `panic` during execution

#HSLIDE

## Go + C?

- Go plugin used in C
  - Already done with -buildmode=shared <!-- .element: class="fragment" -->
- C `.so` loaded in Go
  - For now, error "no plugin module data" <!-- .element: class="fragment" -->

#HSLIDE

## Thank you

Questions?

Find the code and slides on Github: https://github.com/GildasCh/go-plugins-demo
