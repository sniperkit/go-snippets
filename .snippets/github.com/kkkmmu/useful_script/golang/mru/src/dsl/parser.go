package dsl

import (
	"command"
	"fmt"
	"github.com/Workiva/go-datastructures/trie/ctrie"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Parser struct {
	Name string
	*ctrie.Ctrie
	Switch Switch
}

var ParserV8 = &Parser{
	Name:   V8.Name,
	Ctrie:  ctrie.New(nil),
	Switch: V8,
}

var ParserV5 = &Parser{
	Name:   V5.Name,
	Ctrie:  ctrie.New(nil),
	Switch: V5,
}

var ParserM2 = &Parser{
	Name:   M2.Name,
	Ctrie:  ctrie.New(nil),
	Switch: M2,
}

func (p *Parser) Parse(cmd *command.Command) ([]*command.Command, error) {
	res := make([]*command.Command, 0, 1)

	if cmd.CMD == "" {
		return nil, fmt.Errorf("Cannot parser command %v due to empty CMD\n", cmd)
	}

	//Not DSL intruction set(Directy device command)"
	if !strings.HasPrefix(cmd.CMD, "$") {
		c := cmd.Dup()
		c.Delay = 0
		res = append(res, c)
		return res, nil
	}

	return p.parse(cmd)
}

//ospf6(#i).defaultoriginate().always().metric(#metric).routemap(#name)
func (p *Parser) parse(cmd *command.Command) ([]*command.Command, error) {
	res := make([]*command.Command, 0, 1)

	if cmd.CMD == "" {
		return nil, fmt.Errorf("Cannot parser command %v due to empty CMD\n", cmd)
	}

	in := cmd.CMD[1:]

	if i, ok := p.Ctrie.Lookup(MapInstructionToKey(in)); ok {
		if handle, ok := i.(InstructionHandler); ok {
			res := handle(p.Switch, GetInstructionParamsFromCMD(in))
			if len(res) == 0 {
				return nil, fmt.Errorf("Get no command for instruction:%s!\n", cmd.CMD)
			}

			return res, nil
		}
		log.Printf("Invalid result when Get node from DB for instruction: %s\n", cmd.CMD)
	}
	return res, nil
}

func MapInstructionToKey(ins string) []byte {
	var reg = regexp.MustCompile("\\((?P<i>[[:word:][:space:]#$^&.,-_/\\=%]*)\\)")
	matches := reg.FindAllStringSubmatch(ins, -1)
	for _, m := range matches {
		//Here just replace the first one, this is necessary.
		ins = strings.Replace(ins, m[1], "", 1)
	}
	fields := strings.Split(ins, ".")
	var key string
	if len(fields) == 0 {
		key = ins[0:strings.Index(ins, "(")]
	} else {
		for _, f := range fields {
			key += f[0:strings.Index(f, "(")] + "."
		}
	}

	key = strings.TrimRight(key, ".")

	return []byte(key)
}

func GetInstructionParamsFromCMD(ins string) map[string]string {
	ins = strings.Replace(ins, "...)", ")", -1)

	key := string(MapInstructionToKey(ins))
	regs := key + "."

	regs = strings.Replace(regs, ".", "\\((?P<i>[[:word:][:space:]#$%^&*-_=/\\,.+]*)\\).", -1)
	regs = strings.TrimRight(regs, ".")

	reg := regexp.MustCompile(regs)
	matches := reg.FindAllStringSubmatch(ins, -1)

	params := make(map[string]string, 1)
	pas := strings.Split(key, ".")

	for _, m := range matches {
		for j, k := range m {
			if j == 0 {
				continue
			}

			params[pas[j-1]] = k
		}
	}

	return params
}

func DoParserInit(p *Parser) {
	p.Init()
}

func generateInstructionGo() {
	file, err := ioutil.ReadFile("asset/instructions/instruction.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	instruct, err := os.Create("instruction.go")
	if err != nil {
		log.Println("cannot Create file: ", "instruction.go", " ", err.Error())
		return
	}
	defer instruct.Close()

	instruct.WriteString("package dsl\n")
	instruct.WriteString("import ( \n")
	instruct.WriteString(" \"command\"\n")
	instruct.WriteString(" )\n")
	instruct.WriteString("type Instruction struct {\n")
	instruct.WriteString("Name string\n")
	instruct.WriteString("Switch Switch\n")
	instruct.WriteString("}\n")
	for _, l := range lines {
		if l != "" {
			key := string(MapInstructionToKey(l))

			pas := strings.Split(key, ".")
			var fun string
			var parameters string
			for _, p := range pas {
				fun += p
				if p == "no" || p == "No" || p == "NO" || p == "" {
					continue
				}
				parameters += fmt.Sprintf("in[\"%s\"], ", p)
			}

			parameters = strings.TrimRight(parameters, ", ")

			instruct.WriteString("func " + fun + "(sw Switch, in map[string]string) []*command.Command{\n")
			instruct.WriteString(fmt.Sprintf("if len(in) != %d {\n", len(pas)))
			instruct.WriteString("      return nil \n")
			instruct.WriteString("   }\n")
			for _, p := range pas {
				if p == "" {
					continue
				}
				instruct.WriteString(fmt.Sprintf("if _, ok := in[\"%s\"]; !ok {\n", p))
				instruct.WriteString("      return nil \n")
				instruct.WriteString("   }\n")
				if p == "no" {
					continue
				}
			}
			instruct.WriteString(fmt.Sprintf("      return sw.%s(%s)\n", fun, parameters))
			instruct.WriteString("}\n\n")
		}
	}
}

func generateSwitchGo() {
	file, err := ioutil.ReadFile("asset/instructions/instruction.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	sw, err := os.Create("switch.go")
	if err != nil {
		log.Println("cannot Create file: ", "instruction.go", " ", err.Error())
		return
	}

	sw.WriteString("package dsl\n")
	sw.WriteString("import ( \n")
	sw.WriteString(" \"command\"\n")
	sw.WriteString(" )\n")
	sw.WriteString("type Switch interface {\n")

	for _, l := range lines {
		if l != "" {
			key := string(MapInstructionToKey(l))
			pas := strings.Split(key, ".")
			var fun string
			var swparameters string
			for _, p := range pas {
				fun += p
				if p == "no" || p == "No" || p == "NO" || p == "" {
					continue
				}
				swparameters += fmt.Sprintf("%s, ", p)
			}

			swparameters = strings.TrimRight(swparameters, ", ")
			swparameters += " string"
			sw.WriteString(fmt.Sprintf("%s(%s) []*command.Command\n", fun, swparameters))
		}
	}
	sw.WriteString("}\n")
}

func generateSampleGo() {
	file, err := ioutil.ReadFile("asset/instructions/instruction.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	sam, err := os.Create("sample.go")
	if err != nil {
		log.Println("cannot Create file: ", "instruction.go", " ", err.Error())
		return
	}
	sam.WriteString("package dsl\n")
	sam.WriteString("import ( \n")
	sam.WriteString(" \"command\"\n")
	sam.WriteString(" )\n")
	sam.WriteString("type Sample struct {\n")
	sam.WriteString("Name string\n")
	sam.WriteString("}\n")
	sam.WriteString("var SampleDefault = Sample{}\n")

	for _, l := range lines {
		if l != "" {
			key := string(MapInstructionToKey(l))

			pas := strings.Split(key, ".")
			var fun string
			var swparameters string
			for _, p := range pas {
				fun += p
				if p == "no" || p == "No" || p == "NO" || p == "" {
					continue
				}
				swparameters += fmt.Sprintf("%s, ", p)
			}

			swparameters = strings.TrimRight(swparameters, ", ")
			swparameters += " string"
			sam.WriteString(fmt.Sprintf("func (sa Sample) %s(%s) []*command.Command{\n return nil \n}\n\n", fun, swparameters))
		}
	}
}

func generateRegisterGo() {
	file, err := ioutil.ReadFile("asset/instructions/instruction.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	register, err := os.Create("register.go")
	if err != nil {
		log.Println("cannot Create file: ", "register.go", " ", err.Error())
		return
	}

	register.WriteString("package dsl\n")
	register.WriteString("import ( \n")
	register.WriteString(" \"command\"\n")
	register.WriteString(" )\n")

	register.WriteString("type InstructionHandler func(Switch, map[string]string) []*command.Command\n")
	register.WriteString("func (p *Parser) Init() {\n")
	for _, l := range lines {
		if l != "" {
			key := string(MapInstructionToKey(l))
			pas := strings.Split(key, ".")
			var fun string
			for _, p := range pas {
				fun += p
			}

			register.WriteString(fmt.Sprintf("p.Ctrie.Insert([]byte(\"%s\"), InstructionHandler(%s))\n", string(MapInstructionToKey(l)), fun))
		}
	}
	register.WriteString("}\n")
}

func DoGenerate() {
	generateInstructionGo()
	generateSwitchGo()
	generateSampleGo()
	generateRegisterGo()
}

func init() {
	DoGenerate()
	DoParserInit(ParserV8)
	DoParserInit(ParserV5)
	DoParserInit(ParserM2)
}
