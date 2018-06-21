package dsl

import (
	"command"
	"fmt"
)

type DSL struct {
	Parser map[string]*Parser
}

var Engine = DSL{}

func (dsl *DSL) RegisterParser(p *Parser) error {
	if len(dsl.Parser) == 0 {
		dsl.Parser = make(map[string]*Parser, 1)
	}

	if _, ok := dsl.Parser[p.Name]; ok {
		return fmt.Errorf("Same parser: %s already exist!", p.Name)
	}

	dsl.Parser[p.Name] = p

	return nil
}

func (dsl *DSL) ReplaceParser(p *Parser) error {
	if len(dsl.Parser) == 0 {
		dsl.Parser = make(map[string]*Parser, 1)
	}

	dsl.Parser[p.Name] = p

	return nil
}

func (dsl *DSL) DeleteParser(p *Parser) error {
	if len(dsl.Parser) == 0 {
		return fmt.Errorf("Cannot delete parser: %s, due to Empty DB!", p.Name)
	}

	if _, ok := dsl.Parser[p.Name]; !ok {
		return fmt.Errorf("Cannot delete parser: %s, due to Not exist!", p.Name)
	}

	delete(dsl.Parser, p.Name)
	return nil
}

func (dsl *DSL) Parse(device string, cmd *command.Command) ([]*command.Command, error) {
	if dsl.Parser == nil {
		return nil, fmt.Errorf("Cannot parser command due to empty PARSER DB")
	}

	cmd.Validate()
	if _, ok := dsl.Parser[device]; !ok {
		return nil, fmt.Errorf("Cannot parser command due to Unknown command set for device: %s", device)
	}

	return dsl.Parser[device].Parse(cmd)
}

func init() {
	Engine.RegisterParser(ParserV8)
	Engine.RegisterParser(ParserV5)
}
