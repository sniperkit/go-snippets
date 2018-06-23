package main

import (
	//	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Func list</title>
</head>

<body>
{{range .Inst}}
<table border="1">
<p color="red"><b><a href={{.Root}}-----{{.Index}}>{{.Root}}</a></b></p>
{{range .Funcs}}
<p> <tr><a href={{.Name}}-----{{.Index}}>{{.Name}}</a> </tr> </p>
{{end}}
</table>
{{end}}
</body>
</html>
`

var FetchRootName = regexp.MustCompile(`Tracing[[:space:]]+\"(?P<func>[[:word:]_]+)\".*end\.`)
var FetchFunctionName = regexp.MustCompile(`(?P<func>[[:word:]_\.]+\(\))`)
var RemoveCComment = regexp.MustCompile(`(?P<cc>\/\\*.*\*\/)`)
var RemoveCPPComment = regexp.MustCompile(`(?P<cc>\/\/.*)`)

var FetchFunctionCall = regexp.MustCompile(`(?P<func>[[:word:]_\.]+\(\))`)
var RemoveFileHead = FetchRootName
var RemoveFileTail = regexp.MustCompile(`(?P<fe>Ending.*)`)
var RemoveUselessInformation = regexp.MustCompile(`(?P<useless>0[[:space:][:word:]\.\!\\+#)]*\|)`)

type CFunc struct {
	name    string
	body    string
	calling []*CFunc
}

type Temp struct {
	Name  string
	Index int
}
type CInstance struct {
	Index int
	Root  string
	Funcs []*Temp
}

type CInstances struct {
	Inst []*CInstance
}

type CFuncInstance struct {
	index int
	root  *CFunc
	list  map[string]*CFunc
}

type CFuncCallTree struct {
	name      string
	instances []*CFuncInstance
}

/*
    x,y        w ump
	----------------------
        |                    |
       h|                    |
        |                    |
	----------------------
	         dmp
*/
type CFuncGraph struct {
	x, y                   int
	tx, ty                 int
	w, h                   int /* width, height */
	umpx, umpy, dmpx, dmpy int /* up middler port, down middler point */
	name                   string
	level                  int
}

type Graph struct {
	gw, gh         int /* graph size */
	ew, eh         int /* element size */
	sw, sh         int /* Start point */
	deltaw, deltah int /* delta value between to elements */
	pix            int /* pixls */
	rows           map[int]*RowGraph
	graphs         map[string]*CFuncGraph
}

type RowGraph struct {
	x, y    int
	element []*CFuncGraph
}

type CallTreeSPF struct {
	root   string
	vertex map[string]*SPFVertex
	rows   map[int][]*SPFVertex
	level  int
	index  int
}

type SPFVertex struct {
	fn       *CFunc
	level    int
	drawed   bool
	children []*SPFVertex
}

var CallForrest = make([]*CallTreeSPF, 0, 10)

var constGraph = Graph{gw: 4000, gh: 4000, ew: 0, eh: 15, sw: 1000, sh: 40, deltaw: 20, deltah: 40, pix: 5}

var CTreeRoot = CFuncCallTree{}

func (spf *CallTreeSPF) getSubVertex(v *CFunc) {
	if v == nil || len(v.calling) == 0 {
		return
	}

	spf.vertex[v.name].children = make([]*SPFVertex, 0, len(v.calling))
	for _, fi := range v.calling {
		if _, ok := spf.vertex[fi.name]; ok {
			continue
		}

		if fi != nil {
			nv := new(SPFVertex)
			nv.fn = fi
			nv.level = spf.vertex[v.name].level + 1
			if spf.level < nv.level {
				spf.level = nv.level
			}
			nv.drawed = false
			spf.vertex[fi.name] = nv
			spf.vertex[v.name].children = append(spf.vertex[v.name].children, nv)
			spf.getSubVertex(fi)
		}
	}
}

func (spf *CallTreeSPF) getVertex(inst *CFuncInstance) error {
	if r, ok := inst.list[spf.root]; ok {
		spf.vertex = make(map[string]*SPFVertex, 200)
		nv := new(SPFVertex)
		nv.fn = r
		nv.level = 0
		nv.drawed = false
		if spf.level < nv.level {
			spf.level = nv.level
		}
		spf.vertex[spf.root] = nv
		if len(r.calling) != 0 {
			spf.getSubVertex(r)
		}
	}

	spf.rows = make(map[int][]*SPFVertex, spf.level)
	spf.rows[0] = append(spf.rows[0], spf.vertex[spf.root])
	spf.vertex[spf.root].drawed = true
	for l := 1; l <= spf.level; l++ {
		for _, p := range spf.rows[l-1] {
			for _, c := range p.children {
				if spf.vertex[c.fn.name].drawed == true {
					continue
				}
				spf.rows[l] = append(spf.rows[l], c)
				spf.vertex[c.fn.name].drawed = true
			}
		}
	}
	/*
		for _, v := range spf.vertex {
			spf.rows[v.level] = append(spf.rows[v.level], v)
		}
	*/
	return nil
}

func showFuncList(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("Path: ", req.URL.Path)
	if req.URL.Path == "/" {
		instances := CInstances{Inst: make([]*CInstance, 0, len(CTreeRoot.instances))}
		fmt.Println(len(CTreeRoot.instances))
		for _, i := range CTreeRoot.instances {
			inst := new(CInstance)
			inst.Root = i.root.name
			inst.Index = i.index
			inst.Funcs = make([]*Temp, 0, len(i.list))
			for fn, fi := range i.list {
				if len(fi.calling) == 0 {
					continue
				}
				t := new(Temp)
				t.Name = fn
				t.Index = inst.Index
				inst.Funcs = append(inst.Funcs, t)
			}
			instances.Inst = append(instances.Inst, inst)
		}

		/*
			fmt.Println(len(instances.Inst))
			for _, j := range instances.Inst {
				fmt.Println(j.Root + "===============")
				for _, s := range j.Funcs {
					fmt.Println(s)
				}
			}
		*/
		t, err := template.New("main").Parse(tmpl)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, instances)
		return
	}

	createFuncGraph(w, req)
}

func createFuncGraph(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path[1:])

	p := strings.Split(r.URL.Path, "-----")
	index, _ := strconv.ParseInt(p[1], 10, 32)
	fmt.Println(index)
	fmt.Println(p[0])
	if int(index) < len(CTreeRoot.instances) {
		inst := CTreeRoot.instances[index]
		if fn, ok := inst.list[p[0][1:len(p[0])]]; ok {
			spf := new(CallTreeSPF)
			spf.root = fn.name
			spf.level = 0
			spf.index = inst.index
			spf.getVertex(inst)
			buildFuncCallTreeGraph(spf)
			drawFuncCallTree(w, spf)
		} else {
			io.WriteString(w, "Not exist in global list\n")
		}
	} else {
		io.WriteString(w, "----------Not exist in global list--------------\n")
	}
}

func buildRowGraph(level int, row []*SPFVertex) *RowGraph {
	rg := new(RowGraph)
	rlen := 0
	for _, v := range row {
		rlen = rlen + len(v.fn.name)*constGraph.pix + constGraph.deltaw
	}

	rg.x = constGraph.sw - rlen/2
	rg.y = constGraph.sh + (constGraph.eh+constGraph.deltah)*level
	rg.element = make([]*CFuncGraph, 0, len(row))
	sw := rg.x
	sh := rg.y
	for _, v := range row {
		var graph = new(CFuncGraph)
		graph.name = v.fn.name
		graph.x = sw
		graph.y = sh
		graph.w = len(graph.name) * constGraph.pix
		graph.h = constGraph.eh
		graph.tx = graph.x + len(graph.name)*constGraph.pix - (len(graph.name)*constGraph.pix)/2
		graph.ty = graph.y + constGraph.eh - constGraph.pix
		graph.umpx = graph.x + graph.w/2
		graph.umpy = graph.y
		graph.dmpx = graph.x + graph.w/2
		graph.dmpy = graph.y + constGraph.eh
		graph.level = level
		rg.element = append(rg.element, graph)
		constGraph.graphs[graph.name] = graph
		sw += graph.w + constGraph.deltaw
	}

	return rg
}

func buildFuncCallTreeGraph(spf *CallTreeSPF) {
	constGraph.rows = make(map[int]*RowGraph, spf.level)
	constGraph.graphs = make(map[string]*CFuncGraph, len(spf.vertex))
	for i, row := range spf.rows {
		constGraph.rows[i] = buildRowGraph(i, row)
	}
}

func drawFuncCallTree(w http.ResponseWriter, spf *CallTreeSPF) {
	if spf == nil {
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(constGraph.gw, constGraph.gh)
	for _, row := range constGraph.rows {
		for _, g := range row.element {
			s.Roundrect(g.x, g.y, g.w, g.h, 1, 1, "fill:none;stroke:black")
			href := fmt.Sprintf("%s-----%d", g.name, spf.index)
			s.Link(href, g.name)
			s.Text(g.tx, g.ty, g.name, "text-anchor:middle;font-size:5px;fill:black")
			s.LinkEnd()
		}
	}

	for _, v := range spf.vertex {
		for _, subv := range v.children {
			s.Line(constGraph.graphs[v.fn.name].dmpx, constGraph.graphs[v.fn.name].dmpy, constGraph.graphs[subv.fn.name].umpx, constGraph.graphs[subv.fn.name].umpy, "fill:none;stroke:green")
		}
	}
	s.End()
}

func getFuncBody(buf string, figure string) string {
	index := strings.Index(buf, figure)
	body := buf[index+len(figure):]
	var brace_count = 0
	var bsize = 0
	var bstart = 0

	for _, c := range body {
		bsize++
		if bsize == 1 && c != '{' {
			break
		}
		if c == '{' {
			if brace_count == 0 {
				bstart = 1
			}
			brace_count++
		} else if c == '}' {
			brace_count--
			if brace_count == 0 {
				break
			}
		}
	}

	if bsize == 1 {
		body = ""
	} else {
		body = body[bstart : bstart+bsize-2]
	}
	return body
}

func buildCallTree(inst *CFuncInstance) {
	for fn, fi := range inst.list {
		if fi.body == "" {
			inst.list[fn].calling = nil
			continue
		}

		inst.list[fn].calling = make([]*CFunc, 0, 10)
		//fmt.Println(fn + " Body: " + fi.body)
		current := 0
		fstart := current
		insub := 0
		bcount := 0
		for _, c := range fi.body {
			current++
			if c == ')' && insub == 0 {
				fn_name := fi.body[fstart : current-2]
				//				fmt.Println(fn_name)
				inst.list[fn].calling = append(inst.list[fn].calling, inst.list[fn_name])

				if fi.body[current] == '{' {
					insub = 1
					//					bcount++
				}
				continue
			}

			if c == ';' && insub == 0 {
				fstart = current
				continue
			}

			if c == '{' {
				bcount++
				if insub == 0 {
					insub = 1
				}
				continue
			}

			if c == '}' {
				bcount--
				if bcount == 0 {
					insub = 0
					fstart = current
				}
			}
		}
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please input file name")
		os.Exit(-1)
	}

	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		os.Exit(-1)
	}

	buf, _ := ioutil.ReadAll(file)
	//fmt.Println(string(buf))
	//fmt.Println(strings.TrimRight(string(buf), " "))
	s := string(buf)

	s = RemoveCComment.ReplaceAllString(s, "")
	s = RemoveCPPComment.ReplaceAllString(s, "")
	s = RemoveUselessInformation.ReplaceAllString(s, "")

	if FetchRootName.MatchString(s) {
		CTreeRoot.name = FetchRootName.FindStringSubmatch(s)[1][:len(FetchFunctionName.FindStringSubmatch(s)[1])-2]
		CTreeRoot.instances = make([]*CFuncInstance, 0, strings.Count(s, CTreeRoot.name+"()"))
		//fmt.Println(strings.Count(s, CTreeRoot.name+"()"))

		s = RemoveFileTail.ReplaceAllString(s, "")
		s = RemoveFileHead.ReplaceAllString(s, "")
		index := 0
		for n, c := range strings.Split(s, CTreeRoot.name+"()") {
			if n == 0 {
				continue
			}
			cont := 0
			for _, i := range CTreeRoot.instances {
				if strings.EqualFold(c, i.root.body) {
					cont = 1
				}
			}
			if cont == 1 {
				continue
			}

			fmt.Println(c)
			c = strings.Replace(c, "\t", "", -1)
			c = strings.Replace(c, " ", "", -1)
			c = strings.Replace(c, "\n", "", -1)
			fmt.Println(c)
			inst := new(CFuncInstance)
			root := new(CFunc)
			root.name = CTreeRoot.name
			//fmt.Println(c)
			root.body = getFuncBody(root.name+"()"+c, root.name+"()")
			fmt.Println(root.body)
			inst.index = index
			inst.root = root
			CTreeRoot.instances = append(CTreeRoot.instances, inst)
			index++
		}
	}
	for _, i := range CTreeRoot.instances {
		b := i.root.body

		line := strings.Split(b, "()")
		var count = 0
		for _, l := range line {
			if FetchFunctionName.MatchString(l + "()") {
				count = count + 1
			}
			//fmt.Println(l)
		}

		//fmt.Println(b)
		i.list = make(map[string]*CFunc, count)
		i.list[i.root.name] = i.root
		for _, l := range line {
			if FetchFunctionName.MatchString(l + "()") {
				if _, ok := i.list[FetchFunctionName.FindStringSubmatch(l + "()")[1]]; ok {
					continue
				}
				new_fn := new(CFunc)
				new_fn.name = FetchFunctionName.FindStringSubmatch(l + "()")[1][:len(FetchFunctionName.FindStringSubmatch(l + "()")[1])-2]
				new_fn.body = getFuncBody(b, new_fn.name+"()")
				i.list[new_fn.name] = new_fn
			}
		}

		/*
			for _, in := range CTreeRoot.instances {
				//	fmt.Println("++++++++++++++++++++++++++++++++")
				//	fmt.Println(CTreeRoot.name)
				//for _, fn := range in.list {
					//		fmt.Println("----------------------------------")
					//		fmt.Println(fn.name)
				}
			}
		*/

		buildCallTree(i)
		/*
			for fn, fi := range i.list {
				fmt.Println(fn)
				fmt.Println("----------------------------------")
				for _, f := range fi.calling {
					fmt.Println("+" + f.name)
				}
			}
		*/
	}
	http.Handle("/", http.HandlerFunc(showFuncList))
	err = http.ListenAndServe(":2003", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
